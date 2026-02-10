// Simple solid shape shader using vertex buffers
// Used for retained-mode shapes (PrimitiveShape) that need transforms
// Includes edge anti-aliasing using barycentric coordinates

struct VertexInput {
    @location(0) position: vec3<f32>,
    @location(1) local_pos: vec2<f32>,  // barycentric.xy (z = 1 - x - y)
    @location(2) op_code: f32,
    @location(3) radius: f32,
    @location(4) color: vec4<f32>,
    @location(5) tex_coords: vec2<f32>,
}

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(0) color: vec4<f32>,
    @location(1) barycentric: vec3<f32>,
}

// Colors come in as normalized 0..1 (sRGB-ish). We output directly to the swapchain format.
// If the swapchain is sRGB (common), the GPU will handle encoding; doing a manual pow() here
// makes colors/blending look washed/soft.
fn srgbToLinear(c: vec3<f32>) -> vec3<f32> {
    return c;
}

@vertex
fn vs_main(in: VertexInput) -> VertexOutput {
    var output: VertexOutput;
    output.clip_position = vec4<f32>(in.position, 1.0);
    output.color = vec4<f32>(srgbToLinear(in.color.rgb), in.color.a);
    // Reconstruct full barycentric coordinates from xy (z = 1 - x - y)
    output.barycentric = vec3<f32>(in.local_pos.x, in.local_pos.y, 1.0 - in.local_pos.x - in.local_pos.y);
    return output;
}

@fragment
fn fs_main(
    @location(0) color: vec4<f32>,
    @location(1) barycentric: vec3<f32>,
) -> @location(0) vec4<f32> {
    // Edge anti-aliasing using barycentric coordinates
    // The minimum barycentric coordinate indicates distance to the nearest edge
    let min_bary = min(min(barycentric.x, barycentric.y), barycentric.z);

    // Use fwidth for screen-space anti-aliasing width
    let edge_width = fwidth(min_bary) * 0.5;

    // Smoothstep from edge to interior
    let edge_alpha = smoothstep(0.0, edge_width, min_bary);

    return vec4<f32>(color.rgb, color.a * edge_alpha);
}
