#version 410 core

layout(location=0)in vec3 position;
layout(location=1)in vec2 texCoord;

uniform mat4 transform;

out vec2 TexCoord;

void main(){
  gl_Position=transform*vec4(position,1.);
  TexCoord=texCoord;
}
