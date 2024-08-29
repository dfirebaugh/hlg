struct VertexInput {
    @location(0) position: vec3<f32>, // Position in NDC space
    @location(1) color: vec4<f32>,    // Color for the vertex
};

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>, // Clip space position
    @location(0) color: vec4<f32>,               // Pass the color to the fragment shader
};

@vertex
fn vs_main(input: VertexInput) -> VertexOutput {
    var output: VertexOutput;
    output.clip_position = vec4<f32>(input.position, 1.0); // Convert position to clip space
    output.color = input.color;  // Pass the color through to the fragment shader
    return output;
}

@fragment
fn fs_main(input: VertexOutput) -> @location(0) vec4<f32> {
    return input.color;
}

