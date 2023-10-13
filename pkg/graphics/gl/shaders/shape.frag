#version 410 core
in vec2 TexCoord;

out vec4 color;

uniform sampler2D textureSampler;// Texture sampler

void main(){
  color=texture(textureSampler,TexCoord);// Sample the texture using provided texcoords
  //color = vec4(1.0, 0.0, 0.0, 1.0); // Uncomment to output red color for debugging
}