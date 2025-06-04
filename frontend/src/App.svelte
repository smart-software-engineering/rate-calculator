<script>
  import { onMount } from 'svelte';
  import Calculator from './Calculator.svelte';

  // API connection test
  let apiStatus = 'Loading...';
  let apiData = null;

  // Test API connection on component mount
  onMount(async () => {
    try {
      const response = await fetch('/api/v1/health');
      apiData = await response.json();
      apiStatus = 'Connected ✅';
    } catch (error) {
      apiStatus = 'Failed ❌';
      console.error('API connection failed:', error);
    }
  });
</script>

<main>
  <!-- API Status Header -->
  <div class="api-status-header">
    <div class="api-status-compact">
      <span class="status-label">API Status:</span>
      <span class="status-value" class:connected={apiStatus.includes('✅')} class:failed={apiStatus.includes('❌')}>
        {apiStatus}
      </span>
      {#if apiData}
        <span class="api-timestamp">({new Date(apiData.timestamp).toLocaleTimeString()})</span>
      {/if}
    </div>
  </div>

  <!-- Main Calculator -->
  <Calculator />
</main>

<style>
  main {
    width: 100%;
    min-height: 100vh;
    background: #f8f9fa;
  }

  .api-status-header {
    background: white;
    border-bottom: 1px solid #eee;
    padding: 0.75rem 1rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    position: sticky;
    top: 0;
    z-index: 100;
  }

  .api-status-compact {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.9rem;
    max-width: 1200px;
    margin: 0 auto;
  }

  .status-label {
    font-weight: 500;
    color: #666;
  }

  .status-value {
    font-weight: 600;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    background: #f8f9fa;
    border: 1px solid #ddd;
  }

  .status-value.connected {
    background: #d4edda;
    border-color: #c3e6cb;
    color: #155724;
  }

  .status-value.failed {
    background: #f8d7da;
    border-color: #f5c6cb;
    color: #721c24;
  }

  .api-timestamp {
    color: #666;
    font-size: 0.8rem;
    font-style: italic;
  }

  @media (max-width: 768px) {
    .api-status-compact {
      flex-direction: column;
      align-items: flex-start;
      gap: 0.25rem;
    }
  }
</style> 