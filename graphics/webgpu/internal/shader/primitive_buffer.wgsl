const OP_CODE_CIRCLE: f32 = 0.0;
const OP_CODE_ROUNDED_RECT: f32 = 1.0;
const OP_CODE_TRIANGLE: f32 = 2.0;

fn sdCircle(p: vec2<f32>, r: f32) -> f32 {
    return length(p) - r;
}

fn sdRoundedRect(p: vec2<f32>, size: vec2<f32>, radius: f32) -> f32 {
    let q = abs(p) - size + vec2<f32>(radius);
    return length(max(q, vec2<f32>(0.0, 0.0))) - radius;
}

fn sdEquilateralTriangle(p: vec2<f32>) -> f32 {
    let k = sqrt(3.);
    var q: vec2<f32> = vec2<f32>(abs(p.x) - 1.0, p.y + 1. / k);
    if (q.x + k * q.y > 0.) { q = vec2<f32>(q.x - k * q.y, -k * q.x - q.y) / 2.; }
    q.x = q.x - clamp(q.x, -2., 0.);
    return -length(q) * sign(q.y);
}

fn sdTriangle(p: vec2f, p0: vec2f, p1: vec2f, p2: vec2f) -> f32 {
  let e0 = p1 - p0; let e1 = p2 - p1; let e2 = p0 - p2;
  let v0 = p - p0; let v1 = p - p1; let v2 = p - p2;
  let pq0 = v0 - e0 * clamp(dot(v0, e0) / dot(e0, e0), 0., 1.);
  let pq1 = v1 - e1 * clamp(dot(v1, e1) / dot(e1, e1), 0., 1.);
  let pq2 = v2 - e2 * clamp(dot(v2, e2) / dot(e2, e2), 0., 1.);
  let s = sign(e0.x * e2.y - e0.y * e2.x);
  let d = min(min(vec2f(dot(pq0, pq0), s * (v0.x * e0.y - v0.y * e0.x)),
                  vec2f(dot(pq1, pq1), s * (v1.x * e1.y - v1.y * e1.x))),
                  vec2f(dot(pq2, pq2), s * (v2.x * e2.y - v2.y * e2.x)));
  return -sqrt(d.x) * sign(d.y);
}

struct VertexOutput {
    @builtin(position) clipPosition: vec4<f32>,
    @location(0) local_pos: vec2<f32>,
    @location(1) op_code: f32,
    @location(2) radius: f32,
    @location(3) color: vec4<f32>,
}

@vertex
fn vs_main(
    @location(0) in_pos: vec3<f32>,
    @location(1) in_local_pos: vec2<f32>,
    @location(2) in_op_code: f32,
    @location(3) in_radius: f32,
    @location(4) in_color: vec4<f32>
) -> VertexOutput {
    var output: VertexOutput;

    output.clipPosition = vec4<f32>(in_pos, 1.0);

    output.local_pos = in_local_pos;
    output.op_code = in_op_code;
    output.radius = in_radius;
    output.color = in_color;
    return output;
}

@fragment
fn fs_main(
    @location(0) local_pos: vec2<f32>,
    @location(1) op_code: f32,
    @location(2) radius: f32,
    @location(3) color: vec4<f32>
) -> @location(0) vec4<f32> {
    var output_color: vec4<f32> = vec4<f32>(0.0, 0.0, 0.0, 0.0);

    let quad_half_size = vec2<f32>(1.0, 1.0);
    let p = local_pos * quad_half_size * radius;

    if op_code == OP_CODE_CIRCLE {
        let sdf = sdCircle(p, radius);
        if sdf < 0.0 {
            output_color = color;
        }
    } else if op_code == OP_CODE_ROUNDED_RECT {
        let rect_size = vec2<f32>(0.4, 0.3);
        let p = local_pos * vec2<f32>(rect_size.x, rect_size.y);
        let sdf = sdRoundedRect(p, rect_size, radius);
        if sdf < 0.0 {
            output_color = color;
        }
    } else if op_code == OP_CODE_TRIANGLE {
        let sdf = sdEquilateralTriangle(p);
        if sdf < 0.0 {
            output_color = color;
        }
    }
    return output_color;
}

