import { cleanup } from '@testing-library/react'
import '@testing-library/jest-dom/vitest'
import { afterEach } from 'vitest'

// RTL solo registra auto-cleanup con globals:true; con globals:false hay que hacerlo manual.
afterEach(() => cleanup())
