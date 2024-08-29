@group(0) @binding(0) var<uniform> rect_pos: vec2<f32>;
@group(0) @binding(1) var<uniform> rect_size: vec2<f32>;
@group(0) @binding(2) var<uniform> rect_color: vec4<f32>;
@group(0) @binding(3) var<uniform> corner_radius: f32;
@group(0) @binding(4) var<uniform> outline_width: f32;

@vertex
fn vs_main(@location(0) in_pos: vec3<f32>) -> @builtin(position) vec4<f32> {
    return vec4<f32>(in_pos, 1.0);
}

fn sdRect(p: vec2<f32>, size: vec2<f32>, radius: f32) -> f32 {
    let d = abs(p) - size;
    if radius > 0.0 {
        let q = d + vec2<f32>(radius);
        return length(max(q, vec2<f32>(0.0, 0.0))) - radius;
    }
    return max(d.x, d.y);
}

@fragment
fn fs_main(@builtin(position) frag_coord: vec4<f32>) -> @location(0) vec4<f32> {
    let p = frag_coord.xy - (rect_pos + rect_size / 2.0);

    let sdf = sdRect(p, rect_size / 2.0, corner_radius);

    if sdf < 0.0 && sdf > -outline_width {
        return rect_color;
    } else {
        return vec4<f32>(0.0, 0.0, 0.0, 0.0); 
    }
}

