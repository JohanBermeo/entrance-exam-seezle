import { fireEvent, render, screen, within } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { afterAll, afterEach, beforeAll, describe, expect, it, vi } from 'vitest'
import App from './App'

const URL = 'http://localhost:8080/v1/calculations'

function okBody(value: number) {
  return {
    results: { result: { id: 'result', value } },
    outputs: ['result'],
    requestId: 'req-1',
    durationMs: 3,
  }
}

const server = setupServer(http.post(URL, () => HttpResponse.json(okBody(8))))

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

function pressKeys(...names: string[]) {
  const keypad = screen.getByRole('group', { name: 'Teclado' })
  for (const name of names) {
    fireEvent.click(within(keypad).getByRole('button', { name }))
  }
}

describe('App', () => {
  it('calcula 5 + 3 = 8 y deja el resultado como expresión', async () => {
    render(<App />)
    pressKeys('5', '+', '3', '=')
    expect(await screen.findByText('8')).toBeInTheDocument()
    expect(screen.getByText('8', { selector: '.display-expression' })).toBeInTheDocument()
    expect(screen.getByText(/req-1/)).toBeInTheDocument()
  })

  it('encadena operaciones desde el resultado', async () => {
    const seen: string[] = []
    const values = [8, 10]
    server.use(
      http.post(URL, async ({ request }) => {
        const body = (await request.json()) as { expression: string }
        seen.push(body.expression)
        return HttpResponse.json(okBody(values[seen.length - 1] ?? 10))
      }),
    )
    render(<App />)
    pressKeys('5', '+', '3', '=')
    expect(await screen.findByText('8', { selector: '.display-expression' })).toBeInTheDocument()
    pressKeys('+', '2', '=')
    expect(await screen.findByText('10', { selector: '.display-value' })).toBeInTheDocument()
    expect(seen).toEqual(['5+3', '8+2'])
  })

  it('muestra el error de dominio del backend', async () => {
    server.use(
      http.post(URL, () =>
        HttpResponse.json({ code: 'division_by_zero', message: 'divisor cero' }, { status: 422 }),
      ),
    )
    render(<App />)
    pressKeys('5', '÷', '0', '=')
    expect(await screen.findByRole('alert')).toHaveTextContent('division_by_zero')
  })

  it('valida en cliente sin llamar al backend', async () => {
    const handler = vi.fn(() => HttpResponse.json(okBody(0)))
    server.use(http.post(URL, handler))
    render(<App />)
    pressKeys('=')
    expect(await screen.findByRole('alert')).toHaveTextContent('Expresión inválida')
    expect(handler).not.toHaveBeenCalled()
  })

  it('reutiliza una expresión del historial', async () => {
    render(<App />)
    pressKeys('2', '+', '2', '=')
    expect(await screen.findByText('8')).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: /2\+2/ }))
    expect(screen.getByText('2+2', { selector: '.display-expression' })).toBeInTheDocument()
  })
})
