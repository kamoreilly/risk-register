import { afterEach } from 'vitest';
import { cleanup } from '@testing-library/react';
import * as matchers from '@testing-library/jest-dom/matchers';
import { expect } from 'vitest';

// Extend Vitest's expect with jest-dom matchers
expect.extend(matchers);

// Clean up after each test case (e.g., clearing jsdom)
afterEach(() => {
  cleanup();
});