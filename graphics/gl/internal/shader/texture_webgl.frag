#version 300 es

precision highp float;

in vec4 v_color;
in vec2 v_tex_coords;

uniform sampler2D u_texture;

out vec4 frag_color;

void main() {
    vec4 tex_color = texture(u_texture, v_tex_coords);
    frag_color = tex_color * v_color;
}
