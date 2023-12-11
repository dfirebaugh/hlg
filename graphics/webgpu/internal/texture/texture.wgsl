// Vertex shader

struct VertexInput {
    @location(0) position: vec3<f32>,
    @location(1) tex_coords: vec2<f32>,
}

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(0) tex_coords: vec2<f32>,
}

@group(0) @binding(2) var<uniform> u_transform: mat4x4<f32>;

@group(0) @binding(3) var<uniform> u_flipInfo: vec2<f32>;

@group(0) @binding(4) var<uniform> u_clipRect: vec4<f32>;

@vertex
fn vs_main(model: VertexInput) -> VertexOutput {
    var out: VertexOutput;
    out.tex_coords = mix(vec2<f32>(u_clipRect.xy), vec2<f32>(u_clipRect.zw), model.tex_coords);

    // Apply flip based on the uniform
    if (u_flipInfo.x > 0.5) { // Horizontal flip
        out.tex_coords.x = 1.0 - out.tex_coords.x;
    }
    if (u_flipInfo.y > 0.5) { // Vertical flip
        out.tex_coords.y = 1.0 - out.tex_coords.y;
    }

    out.clip_position = u_transform * vec4<f32>(model.position, 1.0);
    return out;
}

// Fragment shader
@group(0) @binding(0)
var t_diffuse: texture_2d<f32>;
@group(0)@binding(1)
var s_diffuse: sampler;

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4<f32> {
    return textureSample(t_diffuse, s_diffuse, in.tex_coords);
}
