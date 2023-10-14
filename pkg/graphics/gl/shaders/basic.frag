#version 330 core

in vec2 TexCoord;// Received from vertex shader

out vec4 FragColor;// Output color

// Texture sampler
uniform sampler2D ourTexture;

void main(){
  FragColor=texture(ourTexture,TexCoord);
}
