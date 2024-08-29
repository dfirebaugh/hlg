@group(0) @binding(0) var<uniform> circle_pos: vec2<f32>;
@group(0) @binding(1) var<uniform> circle_radius: f32;
@group(0) @binding(2) var<uniform> circle_color: vec4<f32>;
@group(0) @binding(3) var<uniform> outline_thickness: f32;

@vertex
fn vs_main(@location(0) in_pos: vec3<f32>) -> @builtin(position) vec4<f32> {
    return vec4<f32>(in_pos, 1.0);
}

fn sdCircle(p: vec2<f32>, radius: f32) -> f32 {
    return length(p) - radius;
}

@fragment
fn fs_main(@builtin(position) frag_coord: vec4<f32>) -> @location(0) vec4<f32> {
    let p = frag_coord.xy - vec2<f32>(circle_pos); 

    let sdf = sdCircle(p, circle_radius);
    if sdf > -outline_thickness && sdf <= 0.0 {
        return circle_color;
    } else {
        return vec4<f32>(0.0, 0.0, 0.0, 0.0); 
    }
}
