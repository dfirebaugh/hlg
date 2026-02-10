#version 410 core

// Primitive buffer fragment shader - SDF rendering for shapes and MSDF text

const float OP_CODE_CIRCLE = 0.0;
const float OP_CODE_ROUNDED_RECT = 1.0;
const float OP_CODE_TRIANGLE = 2.0;
const float OP_CODE_MSDF = 3.0;
const float OP_CODE_SOLID = 4.0;
const float OP_CODE_LINE = 5.0;

in vec2 v_local_pos;
in float v_op_code;
in float v_radius;
in vec4 v_color;
in vec2 v_tex_coords;
in vec2 v_half_size;

uniform sampler2D u_msdf_atlas;
uniform vec4 u_msdf_params; // x=px_range, y=tex_width, z=tex_height, w=msdf_mode

out vec4 frag_color;

// Median of three values (for MSDF)
float median3(float r, float g, float b) {
    return max(min(r, g), min(max(r, g), b));
}

float saturate(float x) {
    return clamp(x, 0.0, 1.0);
}

// SDF functions
float sdCircle(vec2 p, float r) {
    return length(p) - r;
}

float sdRoundedRect(vec2 p, vec2 size, float radius) {
    vec2 q = abs(p) - size + vec2(radius);
    return length(max(q, vec2(0.0))) + min(max(q.x, q.y), 0.0) - radius;
}

float sdEquilateralTriangle(vec2 p) {
    float k = sqrt(3.0);
    vec2 q = vec2(abs(p.x) - 1.0, p.y + 1.0 / k);
    if (q.x + k * q.y > 0.0) {
        q = vec2(q.x - k * q.y, -k * q.x - q.y) / 2.0;
    }
    q.x = q.x - clamp(q.x, -2.0, 0.0);
    return -length(q) * sign(q.y);
}

float sdSegment(vec2 p, vec2 a, vec2 b) {
    vec2 pa = p - a;
    vec2 ba = b - a;
    float h = clamp(dot(pa, ba) / dot(ba, ba), 0.0, 1.0);
    return length(pa - ba * h);
}

float screenPxRange(vec2 tex_coords) {
    float px_range = u_msdf_params.x;
    vec2 tex_size = vec2(u_msdf_params.y, u_msdf_params.z);

    vec2 unit_range = vec2(px_range) / tex_size;
    vec2 screen_tex_size = vec2(1.0) / fwidth(tex_coords);
    // Increased minimum from 1.0 to 1.5 for better AA on low-DPI displays
    return max(0.5 * dot(unit_range, screen_tex_size), 1.5);
}

// Sample MSDF and return signed distance
float sampleMSDF(vec2 uv) {
    vec4 mtsdf = texture(u_msdf_atlas, uv);
    float sd_rgb = median3(mtsdf.r, mtsdf.g, mtsdf.b);
    float sd_a = mtsdf.a;
    return max(sd_rgb, sd_a);
}

void main() {
    frag_color = vec4(0.0, 0.0, 0.0, 0.0);

    int op_code = int(v_op_code + 0.5);

    if (op_code == int(OP_CODE_CIRCLE)) {
        vec2 p = v_local_pos * v_half_size;
        float sdf = sdCircle(p, v_radius);
        float aa = fwidth(sdf) * 0.5;
        float opacity = 1.0 - smoothstep(-aa, aa, sdf);
        if (opacity > 0.005) {
            frag_color = vec4(v_color.rgb, v_color.a * opacity);
        }
    } else if (op_code == int(OP_CODE_ROUNDED_RECT)) {
        vec2 p = v_local_pos * v_half_size;
        float sdf = sdRoundedRect(p, v_half_size, v_radius);
        float aa = fwidth(sdf) * 0.5;
        float opacity = 1.0 - smoothstep(-aa, aa, sdf);
        if (opacity > 0.005) {
            frag_color = vec4(v_color.rgb, v_color.a * opacity);
        }
    } else if (op_code == int(OP_CODE_TRIANGLE)) {
        vec2 p = v_local_pos * v_radius;
        float sdf = sdEquilateralTriangle(p);
        float aa = fwidth(sdf) * 0.5;
        float opacity = 1.0 - smoothstep(-aa, aa, sdf);
        if (opacity > 0.005) {
            frag_color = vec4(v_color.rgb, v_color.a * opacity);
        }
    } else if (op_code == int(OP_CODE_MSDF)) {
        // MTSDF text rendering with 4x supersampling for better quality on low-DPI displays
        float msdf_pxRange = screenPxRange(v_tex_coords);

        // Calculate texel size for supersampling offsets
        vec2 tex_size = vec2(u_msdf_params.y, u_msdf_params.z);
        vec2 texelSize = 1.0 / tex_size;

        // 4x rotated grid supersampling (reduces aliasing better than regular grid)
        // Offsets are ~0.375 texels in a rotated pattern
        vec2 offset = texelSize * 0.375;
        float sd0 = sampleMSDF(v_tex_coords + vec2(-offset.x, -offset.y * 0.5));
        float sd1 = sampleMSDF(v_tex_coords + vec2(offset.x, -offset.y * 0.5));
        float sd2 = sampleMSDF(v_tex_coords + vec2(-offset.x * 0.5, offset.y));
        float sd3 = sampleMSDF(v_tex_coords + vec2(offset.x * 0.5, offset.y));

        // Average the samples
        float sd = (sd0 + sd1 + sd2 + sd3) * 0.25;

        // Convert to screen pixels and apply anti-aliasing
        float screenPxDist = msdf_pxRange * (sd - 0.5);
        float opacity = clamp(screenPxDist + 0.5, 0.0, 1.0);

        if (opacity < 0.005) {
            discard;
        }
        frag_color = vec4(v_color.rgb, v_color.a * opacity);
    } else if (op_code == int(OP_CODE_SOLID)) {
        frag_color = v_color;
    } else if (op_code == int(OP_CODE_LINE)) {
        vec2 p = v_local_pos * v_half_size;
        vec2 line_a = -v_tex_coords;
        vec2 line_b = v_tex_coords;
        float sdf = sdSegment(p, line_a, line_b) - v_radius;
        float aa = fwidth(sdf) * 0.5;
        float opacity = 1.0 - smoothstep(-aa, aa, sdf);
        if (opacity > 0.005) {
            frag_color = vec4(v_color.rgb, v_color.a * opacity);
        }
    }
}
