@group(0) @binding(0) var<uniform> line_start: vec2<f32>;
@group(0) @binding(1) var<uniform> line_end: vec2<f32>;
@group(0) @binding(2) var<uniform> line_color: vec4<f32>;
@group(0) @binding(3) var<uniform> line_thickness: f32;

@vertex
fn vs_main(@location(0) in_pos: vec3<f32>) -> @builtin(position) vec4<f32> {
    return vec4<f32>(in_pos, 1.0);
}

@fragment
fn fs_main(@builtin(position) frag_coord: vec4<f32>) -> @location(0) vec4<f32> {
    let p = frag_coord.xy;
    let ab = line_end - line_start;
    let ap = p - line_start;
    let t = dot(ap, ab) / dot(ab, ab);
    let proj = line_start + ab * clamp(t, 0.0, 1.0);
    let dist = length(proj - p);

    if dist < line_thickness {
        return line_color;
    } else {
        return vec4<f32>(0.0, 0.0, 0.0, 0.0);
    }
}
