import { test, expect } from '@playwright/test';

test.describe('Smoke Test - Vercel Origin CORS', () => {
  test('should allow Vercel origin to access API', async ({ page }) => {
    // Navigate to the health endpoint
    const response = await page.request.get('/api/health');
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data).toHaveProperty('message');
    expect(data.message).toBe('News API is running!');
  });
});