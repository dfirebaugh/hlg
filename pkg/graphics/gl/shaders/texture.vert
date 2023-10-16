#version 330 core

layout(location=0)in vec3 textureAPos;
layout(location=1)in vec3 aColor;
layout(location=2)in vec2 aTexCoord;

out vec2 TexCoord;

uniform vec2 positionOffset;
uniform float rotationAngle;
uniform float scaleWidth;
uniform float scaleHeight;
uniform int windowWidth;
uniform int windowHeight;
uniform float aspectRatioX;
uniform float aspectRatioY;

void main(){
  float cosTheta=cos(rotationAngle);
  float sinTheta=sin(rotationAngle);
  
  vec2 rotatedPos=vec2(
    textureAPos.x*cosTheta-textureAPos.y*sinTheta,
    textureAPos.x*sinTheta+textureAPos.y*cosTheta
  );
  
  vec2 scaledPos=vec2(
    rotatedPos.x*scaleWidth,
    rotatedPos.y*scaleHeight
  );
  
  vec2 aspectAdjustedPos=vec2(
    scaledPos.x*aspectRatioX,
    scaledPos.y*aspectRatioY
  );
  
  vec4 transformedPos=vec4(
    aspectAdjustedPos.x+(2.*(positionOffset.x-.5*scaleWidth*aspectRatioX)/windowWidth-1.),
    aspectAdjustedPos.y-(2.*(positionOffset.y+.5*scaleHeight*aspectRatioY)/windowHeight-1.),
    textureAPos.z,
    1.
  );
  
  gl_Position=transformedPos;
  TexCoord=aTexCoord;
}
