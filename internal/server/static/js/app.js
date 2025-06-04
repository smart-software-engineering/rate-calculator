document.addEventListener('DOMContentLoaded', function() {
    console.log('Rate Calculator loaded - Calculations will be handled server-side with HTMX');
    
    const formInputs = document.querySelectorAll('input, select');
    formInputs.forEach(input => {
        input.addEventListener('change', function() {
            console.log('Form updated - Ready for server-side calculation');
        });
    });

    const percentageSliders = document.querySelectorAll('.percentage-slider');
    percentageSliders.forEach(slider => {
        updateSliderValue(slider);
        
        slider.addEventListener('input', function() {
            updateSliderValue(this);
        });
    });
    
    function updateSliderValue(slider) {
        const valueDisplay = slider.nextElementSibling;
        if (valueDisplay && valueDisplay.classList.contains('slider-value')) {
            valueDisplay.textContent = slider.value + '%';
        }
    }

    const calculateButton = document.getElementById('calculate');
    if (calculateButton) {
        calculateButton.addEventListener('click', function() {
            console.log('Calculate button clicked - Will use HTMX for server-side calculation');
        });
    }
});