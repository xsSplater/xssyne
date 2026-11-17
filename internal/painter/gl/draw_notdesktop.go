//go:build !windows && !linux

package gl

func (p *painter) updateBuffer(vbo Buffer, points []float32) {
	p.ctx.BindBuffer(arrayBuffer, vbo)
	p.logError()
	p.ctx.BufferSubData(arrayBuffer, points)
	p.logError()
}
