<script>
  export let expenseCategories = [];

  function toggleCategory(categoryID) {
    // Toggle visibility of category items
    const categoryElement = document.getElementById(`category-${categoryID}`);
    if (categoryElement) {
      categoryElement.classList.toggle('hidden');
    }
  }

  function handleExpenseChange(categoryIndex, itemIndex, field, event) {
    let value = event.target.value;
    
    if (field === 'Amount') {
      value = parseInt(value) || 0;
    }

    // Update the expense item
    expenseCategories[categoryIndex].Category.Items[itemIndex][field] = value;
    
    // TODO: Will trigger API call to update expenses
    console.log('Expense updated:', field, value);
  }
</script>

<div class="expense-categories">
  {#each expenseCategories as category, categoryIndex}
    <div class="expense-category">
      <div 
        class="category-header" 
        on:click={() => toggleCategory(category.Category.ID)}
        role="button"
        tabindex="0"
        on:keydown={(e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            toggleCategory(category.Category.ID);
          }
        }}
      >
        <h3>{category.Category.Label}</h3>
        <div class="category-summary">
          <div class="summary-item">
            <span class="label">Percentage:</span>
            <span class="value">{category.PercentageTotal}%</span>
          </div>
          <div class="summary-item">
            <span class="label">Yearly:</span>
            <span class="value">€{category.YearlyTotal.toLocaleString()}</span>
          </div>
        </div>
        <span class="toggle-indicator">▼</span>
      </div>
      
      <div class="category-form hidden" id="category-{category.Category.ID}">
        <div class="category-items">
          {#each category.Category.Items as item, itemIndex}
            <div class="expense-item">
              <div class="form-group">
                <label for="expense-{category.Category.ID}-{item.ID}">{item.Label}</label>
                <div class="item-controls">
                  <input 
                    type="number" 
                    id="expense-{category.Category.ID}-{item.ID}"
                    min="0"
                    value={item.Amount}
                    on:input={(e) => handleExpenseChange(categoryIndex, itemIndex, 'Amount', e)}
                  />
                  <select 
                    value={item.Type}
                    on:change={(e) => handleExpenseChange(categoryIndex, itemIndex, 'Type', e)}
                  >
                    <option value="monthly">Monthly</option>
                    <option value="yearly">Yearly</option>
                    <option value="percentage">Percentage</option>
                  </select>
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    </div>
  {/each}
</div>

<style>
  .expense-categories {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .expense-category {
    border: 1px solid #ddd;
    border-radius: 8px;
    overflow: hidden;
  }

  .category-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background: #f8f9fa;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .category-header:hover {
    background: #e9ecef;
  }

  .category-header h3 {
    margin: 0;
    color: #333;
    font-size: 1.1rem;
  }

  .category-summary {
    display: flex;
    gap: 1rem;
  }

  .summary-item {
    display: flex;
    flex-direction: column;
    text-align: center;
  }

  .summary-item .label {
    font-size: 0.8rem;
    color: #666;
  }

  .summary-item .value {
    font-weight: 600;
    color: #333;
  }

  .toggle-indicator {
    font-size: 1.2rem;
    color: #667eea;
    transition: transform 0.2s;
  }

  .category-form {
    padding: 1rem;
    background: white;
    border-top: 1px solid #eee;
  }

  .category-form.hidden {
    display: none;
  }

  .category-items {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .expense-item {
    padding: 0.75rem;
    border: 1px solid #eee;
    border-radius: 4px;
    background: #fafafa;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-group label {
    font-weight: 500;
    color: #333;
  }

  .item-controls {
    display: flex;
    gap: 0.5rem;
  }

  .item-controls input {
    flex: 2;
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
  }

  .item-controls select {
    flex: 1;
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
  }

  .item-controls input:focus,
  .item-controls select:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
  }

  @media (max-width: 768px) {
    .category-header {
      flex-direction: column;
      gap: 0.5rem;
      align-items: stretch;
    }

    .category-summary {
      justify-content: space-around;
    }

    .item-controls {
      flex-direction: column;
    }
  }
</style> 