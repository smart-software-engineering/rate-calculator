<script>
  // Import child components (will create these next)
  import ScheduleForm from './components/ScheduleForm.svelte';
  import ExpenseCategories from './components/ExpenseCategories.svelte';
  import CalculationResult from './components/CalculationResult.svelte';

  // Mock data for now (will be replaced with API calls later)
  let schedules = [
    { ID: '1', Label: 'Standard (40h/week)' },
    { ID: '2', Label: 'Part-time (20h/week)' },
    { ID: '3', Label: 'Freelance (flexible)' }
  ];
  
  let selectedScheduleID = '1';
  let expenseCategories = [
    {
      Category: {
        ID: '1',
        Label: 'Private Expenses',
        Items: [
          { ID: '1', Label: 'Rent', Amount: 1500, Type: 'monthly' },
          { ID: '2', Label: 'Groceries', Amount: 400, Type: 'monthly' }
        ]
      },
      PercentageTotal: 0,
      YearlyTotal: 22800
    },
    {
      Category: {
        ID: '2', 
        Label: 'Professional Expenses',
        Items: [
          { ID: '3', Label: 'Software Licenses', Amount: 100, Type: 'monthly' },
          { ID: '4', Label: 'Insurance', Amount: 2000, Type: 'yearly' }
        ]
      },
      PercentageTotal: 0,
      YearlyTotal: 3200
    }
  ];

  let calculationResult = {
    TotalYearlyExpenses: 26000,
    YearlyTotalWithPercent: 26000,
    HourlyRate: 25
  };

  function handleScheduleChange(event) {
    selectedScheduleID = event.target.value;
    // TODO: Will trigger API call to update schedule
  }

  function handleCalculate() {
    // TODO: Will trigger API call to calculate rates
    console.log('Calculate button clicked');
  }
</script>

<div class="container">
  <header>
    <h1>Freelancer Hourly Rate Calculator</h1>
    <nav>
      <a href="/">← Back to Home</a>
    </nav>
  </header>
  
  <main>
    <div class="main-layout">
      <div class="forms-section">
        <!-- Schedule Section -->
        <div class="schedule-card">
          <div class="section-header">
            <h2>Work Schedule</h2>
            <div class="schedule-selector">
              <select 
                bind:value={selectedScheduleID}
                on:change={handleScheduleChange}
              >
                {#each schedules as schedule}
                  <option value={schedule.ID}>{schedule.Label}</option>
                {/each}
              </select>
            </div>
          </div>
          <ScheduleForm />
        </div>
        
        <!-- Calculator Section -->
        <div class="calculator">
          <div class="calculator-section">
            <h2>Expense Categories</h2>
            <ExpenseCategories {expenseCategories} />
          </div>
          
          <button 
            type="button" 
            class="calculate-button"
            on:click={handleCalculate}
          >
            Calculate Rate
          </button>
        </div>
      </div>
      
      <!-- Results Section -->
      <div class="results-section">
        <div class="result">
          <CalculationResult {calculationResult} />
        </div>
      </div>
    </div>
  </main>
  
  <footer>
    <p>&copy; 2023 Smart Software Engineering SRL</p>
  </footer>
</div>

<style>
  /* Base calculator styles - will be refined */
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 1rem;
  }

  header {
    text-align: center;
    margin-bottom: 2rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #eee;
  }

  header h1 {
    color: #333;
    margin-bottom: 0.5rem;
  }

  nav a {
    color: #667eea;
    text-decoration: none;
  }

  nav a:hover {
    text-decoration: underline;
  }

  .main-layout {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 2rem;
    margin-bottom: 2rem;
  }

  .forms-section {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .schedule-card {
    background: white;
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 1.5rem;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .section-header h2 {
    margin: 0;
    color: #333;
  }

  .schedule-selector select {
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
  }

  .calculator {
    background: white;
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 1.5rem;
  }

  .calculator-section h2 {
    margin-top: 0;
    color: #333;
  }

  .calculate-button {
    background: #667eea;
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: background-color 0.3s;
    margin-top: 1rem;
  }

  .calculate-button:hover {
    background: #5a6fd8;
  }

  .results-section {
    background: white;
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 1.5rem;
  }

  footer {
    text-align: center;
    padding: 2rem 0;
    color: #666;
    border-top: 1px solid #eee;
  }

  @media (max-width: 768px) {
    .main-layout {
      grid-template-columns: 1fr;
    }
    
    .section-header {
      flex-direction: column;
      gap: 1rem;
      align-items: stretch;
    }
  }
</style> 