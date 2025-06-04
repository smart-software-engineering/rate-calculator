document.addEventListener('DOMContentLoaded', function() {
    const amountInput = document.getElementById('amount');
    const rateInput = document.getElementById('rate');
    const termInput = document.getElementById('term');
    const calculateButton = document.getElementById('calculate');
    const resultDiv = document.getElementById('result');

    calculateButton.addEventListener('click', function() {
        // Get values from inputs
        const amount = parseFloat(amountInput.value);
        const rate = parseFloat(rateInput.value);
        const term = parseInt(termInput.value);
        
        // Validate inputs
        if (isNaN(amount) || isNaN(rate) || isNaN(term)) {
            resultDiv.innerHTML = '<p class="error">Please enter valid numbers for all fields</p>';
            return;
        }
        
        if (amount <= 0 || rate <= 0 || term <= 0) {
            resultDiv.innerHTML = '<p class="error">All values must be greater than zero</p>';
            return;
        }
        
        // Calculate monthly rate
        const monthlyRate = rate / 100 / 12;
        
        // Calculate monthly payment
        const monthlyPayment = amount * monthlyRate * Math.pow(1 + monthlyRate, term) / 
                              (Math.pow(1 + monthlyRate, term) - 1);
        
        // Calculate total payment
        const totalPayment = monthlyPayment * term;
        
        // Calculate total interest
        const totalInterest = totalPayment - amount;
        
        // Display results
        resultDiv.innerHTML = `
            <div class="result-item">
                <h3>Monthly Payment</h3>
                <p>$${monthlyPayment.toFixed(2)}</p>
            </div>
            <div class="result-item">
                <h3>Total Payment</h3>
                <p>$${totalPayment.toFixed(2)}</p>
            </div>
            <div class="result-item">
                <h3>Total Interest</h3>
                <p>$${totalInterest.toFixed(2)}</p>
            </div>
        `;
    });
});