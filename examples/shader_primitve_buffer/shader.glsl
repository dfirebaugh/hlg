#vertex
#version 410 core

layout(location = 0) in vec3 in_pos;
layout(location = 1) in vec2 in_local_pos;
layout(location = 2) in float in_op_code;
layout(location = 3) in float in_radius;
layout(location = 4) in vec4 in_color;

out vec2 local_pos;
out float op_code;
out float radius;
out vec4 color;

void main() {
    gl_Position = vec4(in_pos, 1.0);
    local_pos = in_local_pos;
    op_code = in_op_code;
    radius = in_radius;
    color = in_color;
}
#fragment
#version 410 core

in vec2 local_pos;
in float op_code;
in float radius;
in vec4 color;

out vec4 fragColor;

const float OP_CODE_CIRCLE = 0.0;
const float OP_CODE_ROUNDED_RECT = 1.0;
const float OP_CODE_TRIANGLE = 2.0;

float sdCircle(vec2 p, float r) {
    return length(p) - r;
}

float sdRoundedRect(vec2 p, vec2 size, float r) {
    vec2 q = abs(p) - size + vec2(r);
    return length(max(q, vec2(0.0, 0.0))) - r;
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

void main() {
    vec4 output_color = vec4(0.0, 0.0, 0.0, 0.0);

    vec2 quad_half_size = vec2(1.0, 1.0);
    vec2 p = local_pos * quad_half_size * radius;

    if (op_code == OP_CODE_CIRCLE) {
        float sdf = sdCircle(p, radius);
        if (sdf < 0.0) {
            output_color = color;
        }
    } else if (op_code == OP_CODE_ROUNDED_RECT) {
        vec2 rect_size = vec2(0.4, 0.3);
        vec2 pp = local_pos * vec2(rect_size.x, rect_size.y);
        float sdf = sdRoundedRect(pp, rect_size, radius);
        if (sdf < 0.0) {
            output_color = color;
        }
    } else if (op_code == OP_CODE_TRIANGLE) {
        float sdf = sdEquilateralTriangle(p);
        if (sdf < 0.0) {
            output_color = color;
        }
    }
    fragColor = output_color;
}
