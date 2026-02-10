#version 410 core

layout(location = 0) in vec3 a_position;
layout(location = 1) in vec4 a_color;
layout(location = 2) in vec2 a_tex_coords;

uniform mat4 u_transform;
uniform vec2 u_flip_info;
uniform vec4 u_clip_rect;

out vec4 v_color;
out vec2 v_tex_coords;

void main() {
    v_tex_coords = a_tex_coords * (u_clip_rect.zw - u_clip_rect.xy) + u_clip_rect.xy;

    if (u_flip_info.x > 0.5) {
        v_tex_coords.x = 1.0 - v_tex_coords.x;
    }
    if (u_flip_info.y > 0.5) {
        v_tex_coords.y = 1.0 - v_tex_coords.y;
    }

    gl_Position = u_transform * vec4(a_position, 1.0);
    v_color = a_color;
}
