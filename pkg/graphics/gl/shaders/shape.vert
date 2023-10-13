#version 410 core
layout(location=0)in vec2 position;
layout(location=1)in vec2 texCoord;// Texture coordinates

uniform mat4 transform;

out vec2 TexCoord;

void main(){
  vec2 flippedPosition=vec2(position.x,-position.y);
  gl_Position=transform*vec4(flippedPosition,0.,1.);
  TexCoord=texCoord;// Pass texture coordinate to fragment shader
}