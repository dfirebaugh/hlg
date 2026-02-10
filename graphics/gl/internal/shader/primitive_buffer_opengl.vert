#version 410 core

// Primitive buffer shader - vertex buffer based for OpenGL 4.1 compatibility
// Uses PrimitiveVertex format: Position[3], LocalPosition[2], OpCode, Radius, Color[4], TexCoords[2], HalfSize[2]

layout(location = 0) in vec3 a_position;
layout(location = 1) in vec2 a_local_pos;
layout(location = 2) in float a_op_code;
layout(location = 3) in float a_radius;
layout(location = 4) in vec4 a_color;
layout(location = 5) in vec2 a_tex_coords;
layout(location = 6) in vec2 a_half_size;

uniform vec2 u_screen_size;

out vec2 v_local_pos;
out float v_op_code;
out float v_radius;
out vec4 v_color;
out vec2 v_tex_coords;
out vec2 v_half_size;

void main() {
    gl_Position = vec4(a_position.xy, 0.0, 1.0);
    v_local_pos = a_local_pos;
    v_op_code = a_op_code;
    v_radius = a_radius;
    v_color = a_color;
    v_tex_coords = a_tex_coords;
    v_half_size = a_half_size;
}
