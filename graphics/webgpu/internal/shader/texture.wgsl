
struct VertexInput {
    @location(0) position: vec3<f32>,
    @location(1) color: vec4<f32>,
    @location(2) tex_coords: vec2<f32>,
};

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(1) color: vec4<f32>,
    @location(2) tex_coords: vec2<f32>,
};

@group(0) @binding(2) var<uniform> u_transform: mat4x4<f32>;

@group(0) @binding(3) var<uniform> u_flipInfo: vec2<f32>;

@group(0) @binding(4) var<uniform> u_clipRect: vec4<f32>;

@group(0) @binding(0) var t_diffuse: texture_2d<f32>;

@group(0) @binding(1) var s_diffuse: sampler;

@vertex
fn vs_main(model: VertexInput) -> VertexOutput {
    var out: VertexOutput;

    out.tex_coords = model.tex_coords * (u_clipRect.zw - u_clipRect.xy) + u_clipRect.xy;

    if (u_flipInfo.x > 0.5) {
        out.tex_coords.x = 1.0 - out.tex_coords.x;
    }
    if (u_flipInfo.y > 0.5) {
        out.tex_coords.y = 1.0 - out.tex_coords.y;
    }

    out.clip_position = u_transform * vec4<f32>(model.position, 1.0);
    out.color = model.color;
    return out;
}

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4<f32> {
    let tex_color = textureSample(t_diffuse, s_diffuse, in.tex_coords);
    return tex_color * in.color;
}

