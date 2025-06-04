<script>
  export let calculationResult = {
    TotalYearlyExpenses: 0,
    YearlyTotalWithPercent: 0,
    HourlyRate: 0
  };

  // Calculate some derived values for display
  $: monthlyExpenses = Math.round(calculationResult.TotalYearlyExpenses / 12);
  $: weeklyExpenses = Math.round(calculationResult.TotalYearlyExpenses / 52);
  $: dailyExpenses = Math.round(calculationResult.TotalYearlyExpenses / 365);
</script>

<div class="calculation-result">
  <h3>📊 Rate Calculation Results</h3>
  
  <div class="result-grid">
    <!-- Primary Result -->
    <div class="result-card primary">
      <h4>💰 Recommended Hourly Rate</h4>
      <div class="rate-display">
        <span class="currency">€</span>
        <span class="amount">{calculationResult.HourlyRate}</span>
        <span class="period">/hour</span>
      </div>
    </div>

    <!-- Expense Breakdown -->
    <div class="result-card">
      <h4>📋 Expense Breakdown</h4>
      <div class="breakdown-list">
        <div class="breakdown-item">
          <span class="label">Yearly Total:</span>
          <span class="value">€{calculationResult.TotalYearlyExpenses.toLocaleString()}</span>
        </div>
        <div class="breakdown-item">
          <span class="label">Monthly:</span>
          <span class="value">€{monthlyExpenses.toLocaleString()}</span>
        </div>
        <div class="breakdown-item">
          <span class="label">Weekly:</span>
          <span class="value">€{weeklyExpenses.toLocaleString()}</span>
        </div>
        <div class="breakdown-item">
          <span class="label">Daily:</span>
          <span class="value">€{dailyExpenses.toLocaleString()}</span>
        </div>
      </div>
    </div>

    <!-- Rate Scenarios -->
    <div class="result-card">
      <h4>🎯 Rate Scenarios</h4>
      <div class="scenario-list">
        <div class="scenario-item">
          <span class="scenario-label">Minimum (Survival):</span>
          <span class="scenario-rate">€{Math.round(calculationResult.HourlyRate * 0.8)}/hr</span>
        </div>
        <div class="scenario-item">
          <span class="scenario-label">Recommended:</span>
          <span class="scenario-rate primary">€{calculationResult.HourlyRate}/hr</span>
        </div>
        <div class="scenario-item">
          <span class="scenario-label">Growth (20% margin):</span>
          <span class="scenario-rate">€{Math.round(calculationResult.HourlyRate * 1.2)}/hr</span>
        </div>
      </div>
    </div>

    <!-- Additional Info -->
    <div class="result-card">
      <h4>ℹ️ Additional Information</h4>
      <div class="info-list">
        <p><strong>🏠 Remote Work:</strong> Use the recommended rate</p>
        <p><strong>🏢 On-site Work:</strong> Add travel time and expenses</p>
        <p><strong>📈 Yearly Review:</strong> Adjust rates based on inflation and skill growth</p>
      </div>
    </div>
  </div>

  <div class="result-actions">
    <button class="action-button secondary">📄 Export PDF</button>
    <button class="action-button secondary">📊 Detailed Report</button>
    <button class="action-button primary">💾 Save Calculation</button>
  </div>
</div>

<style>
  .calculation-result {
    padding: 1rem;
  }

  .calculation-result h3 {
    margin: 0 0 1.5rem 0;
    color: #333;
    text-align: center;
  }

  .result-grid {
    display: grid;
    gap: 1rem;
    margin-bottom: 2rem;
  }

  .result-card {
    background: white;
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .result-card.primary {
    border-color: #667eea;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  .result-card h4 {
    margin: 0 0 1rem 0;
    font-size: 1.1rem;
  }

  .rate-display {
    display: flex;
    align-items: baseline;
    justify-content: center;
    gap: 0.25rem;
  }

  .rate-display .currency {
    font-size: 1.5rem;
    font-weight: 600;
  }

  .rate-display .amount {
    font-size: 3rem;
    font-weight: 700;
  }

  .rate-display .period {
    font-size: 1.2rem;
    opacity: 0.9;
  }

  .breakdown-list,
  .scenario-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .breakdown-item,
  .scenario-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 0;
    border-bottom: 1px solid #f0f0f0;
  }

  .breakdown-item:last-child,
  .scenario-item:last-child {
    border-bottom: none;
  }

  .label,
  .scenario-label {
    color: #666;
    font-weight: 500;
  }

  .value,
  .scenario-rate {
    font-weight: 600;
    color: #333;
  }

  .scenario-rate.primary {
    color: #667eea;
    font-weight: 700;
  }

  .info-list p {
    margin: 0.5rem 0;
    font-size: 0.9rem;
    color: #666;
    line-height: 1.5;
  }

  .result-actions {
    display: flex;
    gap: 0.75rem;
    justify-content: center;
    flex-wrap: wrap;
  }

  .action-button {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 6px;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-button.primary {
    background: #667eea;
    color: white;
  }

  .action-button.primary:hover {
    background: #5a6fd8;
    transform: translateY(-1px);
  }

  .action-button.secondary {
    background: #f8f9fa;
    color: #667eea;
    border: 1px solid #667eea;
  }

  .action-button.secondary:hover {
    background: #667eea;
    color: white;
  }

  @media (max-width: 768px) {
    .result-actions {
      flex-direction: column;
    }

    .action-button {
      width: 100%;
    }

    .rate-display .amount {
      font-size: 2.5rem;
    }
  }
</style> 