package shader

import (
	"log"
	"sync"

	"github.com/dfirebaugh/hlg/graphics"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type ShaderManager struct {
	device      *wgpu.Device
	shaderCache map[graphics.ShaderHandle]*wgpu.ShaderModule
	nextHandle  graphics.ShaderHandle
	cacheMutex  sync.RWMutex
}

func NewShaderManager(device *wgpu.Device) *ShaderManager {
	return &ShaderManager{
		device:      device,
		shaderCache: make(map[graphics.ShaderHandle]*wgpu.ShaderModule),
		nextHandle:  1,
	}
}

func (sm *ShaderManager) CompileShader(shaderCode string) graphics.ShaderHandle {
	if sm == nil {
		log.Fatal("ShaderManager is nil")
	}

	shaderModule, err := sm.device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: shaderCode},
	})
	if err != nil {
		log.Fatalf("Failed to create shader module: %v", err)
	}

	sm.cacheMutex.Lock()
	handle := sm.nextHandle
	sm.nextHandle++
	sm.shaderCache[handle] = shaderModule
	sm.cacheMutex.Unlock()

	return handle
}

func (sm *ShaderManager) GetShader(handle graphics.ShaderHandle) *wgpu.ShaderModule {
	sm.cacheMutex.RLock()
	shader, exists := sm.shaderCache[handle]
	sm.cacheMutex.RUnlock()

	if !exists {
		log.Printf("Shader not found for handle: %d", handle)
		return nil
	}

	return shader
}

func (sm *ShaderManager) ReleaseShaders() {
	for _, s := range sm.shaderCache {
		if s == nil {
			continue
		}
		// s.Release()
		s = nil
	}
}
