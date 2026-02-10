// Primitive buffer shader - One draw call UI rendering
// Based on: https://ruby0x1.github.io/machinery_blog_archive/post/ui-rendering-using-primitive-buffers/
//
// Key insight: Store one compact primitive per shape in a storage buffer.
// The vertex shader constructs 6 vertices (2 triangles) per primitive using vertex_index.
// This reduces memory from 312 bytes/shape to 48 bytes/shape (~6.5x reduction).

const OP_CODE_CIRCLE: f32 = 0.0;
const OP_CODE_ROUNDED_RECT: f32 = 1.0;
const OP_CODE_TRIANGLE: f32 = 2.0;
const OP_CODE_MSDF: f32 = 3.0;
const OP_CODE_SOLID: f32 = 4.0;
const OP_CODE_LINE: f32 = 5.0;

// Compact primitive data (64 bytes per primitive, 16-byte aligned)
// MUST match Go Primitive struct layout exactly!
struct Primitive {
    x: f32,           // offset 0
    y: f32,           // offset 4
    w: f32,           // offset 8
    h: f32,           // offset 12
    color: vec4<f32>, // offset 16 (16-byte aligned)
    radius: f32,      // offset 32
    op_code: f32,     // offset 36
    _pad1: vec2<f32>, // offset 40 (padding to align extra)
    extra: vec4<f32>, // offset 48 (16-byte aligned) - for MSDF (u0, v0, u_size, v_size); for shapes (half_w, half_h, 0, 0)
}

// Storage buffer containing all primitives
@group(0) @binding(0) var<storage, read> primitives: array<Primitive>;

// Uniforms for screen size (needed for NDC conversion)
@group(0) @binding(1) var<uniform> screen_size: vec2<f32>;

// MSDF atlas resources
@group(0) @binding(2) var t_msdf_atlas: texture_2d<f32>;
@group(0) @binding(3) var s_msdf_atlas: sampler;
// x=px_range, y=tex_width, z=tex_height, w=msdf_mode (0=median RGB/MSDF, 1=alpha/true SDF, 2=visualize RGB)
@group(0) @binding(4) var<uniform> u_msdf_params: vec4<f32>;

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(0) local_pos: vec2<f32>,
    @location(1) op_code: f32,
    @location(2) radius: f32,
    @location(3) color: vec4<f32>,
    @location(4) tex_coords: vec2<f32>,
    @location(5) half_size: vec2<f32>,
}

// Colors come in as normalized 0..1 (sRGB-ish). We output directly to the swapchain format.
// If the swapchain is sRGB (common), the GPU will handle encoding; doing a manual pow() here
// makes colors/blending look washed/soft.
fn srgbToLinear(c: vec3<f32>) -> vec3<f32> {
    return c;
}

// Median of three values (used for MSDF reconstruction)
fn median3(r: f32, g: f32, b: f32) -> f32 {
    return max(min(r, g), min(max(r, g), b));
}

fn saturate(x: f32) -> f32 {
    return clamp(x, 0.0, 1.0);
}

// Convert screen coordinates to NDC
fn screenToNDC(pos: vec2<f32>, screen: vec2<f32>) -> vec2<f32> {
    let normalized = pos / screen;
    return vec2<f32>(normalized.x * 2.0 - 1.0, 1.0 - normalized.y * 2.0);
}

// Corner positions for a quad (2 triangles, 6 vertices)
// Tri 1: 0=BL, 1=BR, 2=TL | Tri 2: 3=TL, 4=BR, 5=TR
fn getCornerOffset(corner_index: u32) -> vec2<f32> {
    switch corner_index {
        case 0u: { return vec2<f32>(0.0, 1.0); }  // bottom-left
        case 1u: { return vec2<f32>(1.0, 1.0); }  // bottom-right
        case 2u: { return vec2<f32>(0.0, 0.0); }  // top-left
        case 3u: { return vec2<f32>(0.0, 0.0); }  // top-left
        case 4u: { return vec2<f32>(1.0, 1.0); }  // bottom-right
        case 5u: { return vec2<f32>(1.0, 0.0); }  // top-right
        default: { return vec2<f32>(0.0, 0.0); }
    }
}

// Local position for SDF (-1 to 1 range)
fn getLocalPos(corner_index: u32) -> vec2<f32> {
    switch corner_index {
        case 0u: { return vec2<f32>(-1.0, -1.0); }
        case 1u: { return vec2<f32>(1.0, -1.0); }
        case 2u: { return vec2<f32>(-1.0, 1.0); }
        case 3u: { return vec2<f32>(-1.0, 1.0); }
        case 4u: { return vec2<f32>(1.0, -1.0); }
        case 5u: { return vec2<f32>(1.0, 1.0); }
        default: { return vec2<f32>(0.0, 0.0); }
    }
}

// UV coordinates for MSDF (0 to 1 range, standard texture convention: V=0 at top, V=1 at bottom)
fn getUVOffset(corner_index: u32) -> vec2<f32> {
    switch corner_index {
        case 0u: { return vec2<f32>(0.0, 1.0); }  // bottom-left (V=1, bottom of texture)
        case 1u: { return vec2<f32>(1.0, 1.0); }  // bottom-right (V=1, bottom of texture)
        case 2u: { return vec2<f32>(0.0, 0.0); }  // top-left (V=0, top of texture)
        case 3u: { return vec2<f32>(0.0, 0.0); }  // top-left (V=0, top of texture)
        case 4u: { return vec2<f32>(1.0, 1.0); }  // bottom-right (V=1, bottom of texture)
        case 5u: { return vec2<f32>(1.0, 0.0); }  // top-right (V=0, top of texture)
        default: { return vec2<f32>(0.0, 0.0); }
    }
}

@vertex
fn vs_main(@builtin(vertex_index) vertex_index: u32) -> VertexOutput {
    let prim_index = vertex_index / 6u;
    let corner_index = vertex_index % 6u;

    let prim = primitives[prim_index];
    let corner = getCornerOffset(corner_index);
    let screen_pos = vec2<f32>(prim.x, prim.y) + corner * vec2<f32>(prim.w, prim.h);
    let ndc = screenToNDC(screen_pos, screen_size);
    let local_pos = getLocalPos(corner_index);

    var output: VertexOutput;
    output.clip_position = vec4<f32>(ndc, 0.0, 1.0);
    output.local_pos = local_pos;
    output.op_code = prim.op_code;
    output.radius = prim.radius;
    output.color = vec4<f32>(srgbToLinear(prim.color.rgb), prim.color.a);
    output.half_size = vec2<f32>(prim.w, prim.h) * 0.5;

    // For MSDF: extra stores UV base (xy) and UV size (zw)
    if prim.op_code == OP_CODE_MSDF {
        let uv_offset = getUVOffset(corner_index);

        // Direct UV coordinates without inset (banana-c approach)
        // MSDF with proper pixel range doesn't need UV inset
        let uv_base = prim.extra.xy;
        let uv_size = prim.extra.zw;

        output.tex_coords = uv_base + uv_offset * uv_size;
    } else {
        // For shapes: extra.xy contains half_size (used in fragment shader for SDF)
        output.tex_coords = prim.extra.xy;
    }

    return output;
}

// SDF functions
fn sdCircle(p: vec2<f32>, r: f32) -> f32 {
    return length(p) - r;
}

fn sdRoundedRect(p: vec2<f32>, size: vec2<f32>, radius: f32) -> f32 {
    let q = abs(p) - size + vec2<f32>(radius);
    return length(max(q, vec2<f32>(0.0, 0.0))) + min(max(q.x, q.y), 0.0) - radius;
}

fn sdEquilateralTriangle(p: vec2<f32>) -> f32 {
    let k = sqrt(3.);
    var q: vec2<f32> = vec2<f32>(abs(p.x) - 1.0, p.y + 1. / k);
    if (q.x + k * q.y > 0.) { q = vec2<f32>(q.x - k * q.y, -k * q.x - q.y) / 2.; }
    q.x = q.x - clamp(q.x, -2., 0.);
    return -length(q) * sign(q.y);
}

// SDF for a line segment (capsule shape)
// p = point to test, a = start, b = end
fn sdSegment(p: vec2<f32>, a: vec2<f32>, b: vec2<f32>) -> f32 {
    let pa = p - a;
    let ba = b - a;
    let h = clamp(dot(pa, ba) / dot(ba, ba), 0.0, 1.0);
    return length(pa - ba * h);
}

fn median(r: f32, g: f32, b: f32) -> f32 {
    return max(min(r, g), min(max(r, g), b));
}

fn screenPxRange(tex_coords: vec2<f32>) -> f32 {
    let px_range = u_msdf_params.x;
    let tex_size = vec2<f32>(u_msdf_params.y, u_msdf_params.z);

    // Standard MSDF screenPxRange calculation (msdfgen reference)
    let unit_range = vec2<f32>(px_range) / tex_size;
    let screen_tex_size = vec2<f32>(1.0) / fwidth(tex_coords);
    // Increased minimum from 1.0 to 1.5 for better AA on low-DPI displays
    return max(0.5 * dot(unit_range, screen_tex_size), 1.5);
}

// Sample MSDF and return signed distance
fn sampleMSDF(uv: vec2<f32>) -> f32 {
    let mtsdf = textureSample(t_msdf_atlas, s_msdf_atlas, uv);
    let sd_rgb = median3(mtsdf.r, mtsdf.g, mtsdf.b);
    let sd_a = mtsdf.a;
    return max(sd_rgb, sd_a);
}

@fragment
fn fs_main(
    @location(0) local_pos: vec2<f32>,
    @location(1) op_code: f32,
    @location(2) radius: f32,
    @location(3) color: vec4<f32>,
    @location(4) tex_coords: vec2<f32>,
    @location(5) half_size: vec2<f32>,
) -> @location(0) vec4<f32> {
    var output_color: vec4<f32> = vec4<f32>(0.0, 0.0, 0.0, 0.0);

    // MTSDF rendering with 4x supersampling for better quality on low-DPI displays
    let msdf_pxRange = screenPxRange(tex_coords);

    // Calculate texel size for supersampling offsets
    let tex_size = vec2<f32>(u_msdf_params.y, u_msdf_params.z);
    let texelSize = 1.0 / tex_size;

    // 4x rotated grid supersampling (reduces aliasing better than regular grid)
    let offset = texelSize * 0.375;
    let sd0 = sampleMSDF(tex_coords + vec2<f32>(-offset.x, -offset.y * 0.5));
    let sd1 = sampleMSDF(tex_coords + vec2<f32>(offset.x, -offset.y * 0.5));
    let sd2 = sampleMSDF(tex_coords + vec2<f32>(-offset.x * 0.5, offset.y));
    let sd3 = sampleMSDF(tex_coords + vec2<f32>(offset.x * 0.5, offset.y));

    // Average the samples
    let sd = (sd0 + sd1 + sd2 + sd3) * 0.25;

    // Convert to screen pixels and apply anti-aliasing
    let screenPxDist = msdf_pxRange * (sd - 0.5);
    let msdf_opacity_rgb = clamp(screenPxDist + 0.5, 0.0, 1.0);

    // Alpha-only opacity for fallback mode (single sample is fine for fallback)
    let mtsdf_a = textureSample(t_msdf_atlas, s_msdf_atlas, tex_coords);
    let sd_a = mtsdf_a.a;
    let pxDist_a = msdf_pxRange * (sd_a - 0.5);
    let msdf_opacity_a = clamp(pxDist_a + 0.5, 0.0, 1.0);

    let circle_p = local_pos * half_size;
    let circle_sdf = sdCircle(circle_p, radius);
    let circle_aa = fwidth(circle_sdf) * 0.5;

    let rect_p = local_pos * half_size;
    let rect_sdf = sdRoundedRect(rect_p, half_size, radius);
    let rect_aa = fwidth(rect_sdf) * 0.5;

    let tri_p = local_pos * radius;
    let tri_sdf = sdEquilateralTriangle(tri_p);
    let tri_aa = fwidth(tri_sdf) * 0.5;

    // Line SDF: tex_coords stores (cos*halfLen, sin*halfLen), radius stores half_width
    let line_p = local_pos * half_size;
    // tex_coords contains the line endpoint offset from center in pixels
    let line_a = -tex_coords;
    let line_b = tex_coords;
    let line_sdf = sdSegment(line_p, line_a, line_b) - radius;
    let line_aa = fwidth(line_sdf) * 0.5;

    if op_code == OP_CODE_CIRCLE {
        let opacity = 1.0 - smoothstep(-circle_aa, circle_aa, circle_sdf);
        if opacity > 0.005 {
            output_color = vec4<f32>(color.rgb, color.a * opacity);
        }
    } else if op_code == OP_CODE_ROUNDED_RECT {
        let opacity = 1.0 - smoothstep(-rect_aa, rect_aa, rect_sdf);
        if opacity > 0.005 {
            output_color = vec4<f32>(color.rgb, color.a * opacity);
        }
    } else if op_code == OP_CODE_TRIANGLE {
        let opacity = 1.0 - smoothstep(-tri_aa, tri_aa, tri_sdf);
        if opacity > 0.005 {
            output_color = vec4<f32>(color.rgb, color.a * opacity);
        }
    } else if op_code == OP_CODE_MSDF {
        // MTSDF rendering - uses multi-channel signed distance field for sharp corners
        // Mode 0: median(RGB) - MSDF reconstruction for sharp corners (default)
        // Mode 1: alpha channel only (true SDF fallback)
        // Mode 2: visualize RGB channels directly (for debugging atlas)
        // Mode 3: hard threshold (no AA)
        // Mode 4: enhanced supersampling (8 samples)
        // Mode 5: adaptive AA width
        // Mode 6: sharp mode for small text
        // Mode 7: soft mode for large text
        // Mode 8: crisp mode - ultra tight AA for small text
        let msdf_mode = u_msdf_params.w;

        // Simplified mode handling (banana-c style - single sample is sufficient)
        if msdf_mode >= 2.5 {
            // Mode 3+: Hard threshold test (no AA) - for debugging atlas data
            let hard_opacity = select(0.0, 1.0, sd > 0.5);
            output_color = vec4<f32>(color.rgb, color.a * hard_opacity);
            return output_color;
        } else if msdf_mode >= 1.5 {
            // Mode 2: Visualize RGB directly (for debugging)
            output_color = vec4<f32>(mtsdf.rgb, 1.0);
            return output_color;
        } else if msdf_mode >= 0.5 {
            // Mode 1: alpha-only (true SDF) fallback
            if msdf_opacity_a < 0.005 { discard; }
            output_color = vec4<f32>(color.rgb, color.a * msdf_opacity_a);
            return output_color;
        }

        // Mode 0: standard MSDF (median RGB) - the default banana-c approach
        if msdf_opacity_rgb < 0.005 { discard; }
        output_color = vec4<f32>(color.rgb, color.a * msdf_opacity_rgb);
    } else if op_code == OP_CODE_SOLID {
        output_color = color;
    } else if op_code == OP_CODE_LINE {
        let opacity = 1.0 - smoothstep(-line_aa, line_aa, line_sdf);
        if opacity > 0.005 {
            output_color = vec4<f32>(color.rgb, color.a * opacity);
        }
    }

    return output_color;
}
