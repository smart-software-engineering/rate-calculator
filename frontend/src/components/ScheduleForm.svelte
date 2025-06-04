<script>
  // Mock schedule data (will be replaced with props/API data)
  let schedule = {
    HoursPerWeek: 40,
    VacationDays: 25,
    PublicHolidays: 10,
    EducationDays: 5,
    SickDays: 5
  };

  function handleInputChange(field, event) {
    schedule[field] = field === 'HoursPerWeek' ? 
      parseFloat(event.target.value) : 
      parseInt(event.target.value);
    
    // TODO: Will trigger API call to update schedule
    console.log('Schedule updated:', field, schedule[field]);
  }
</script>

<div class="schedule-form">
  <div class="form-group">
    <label for="hoursPerWeek">Hours per Week:</label>
    <input 
      type="number" 
      id="hoursPerWeek"
      step="0.5"
      min="1" 
      max="80"
      value={schedule.HoursPerWeek}
      on:input={(e) => handleInputChange('HoursPerWeek', e)}
    />
  </div>

  <div class="form-group">
    <label for="vacationDays">Vacation Days per Year:</label>
    <input 
      type="number" 
      id="vacationDays"
      min="0" 
      max="365"
      value={schedule.VacationDays}
      on:input={(e) => handleInputChange('VacationDays', e)}
    />
  </div>

  <div class="form-group">
    <label for="publicHolidays">Public Holidays:</label>
    <input 
      type="number" 
      id="publicHolidays"
      min="0" 
      max="50"
      value={schedule.PublicHolidays}
      on:input={(e) => handleInputChange('PublicHolidays', e)}
    />
  </div>

  <div class="form-group">
    <label for="educationDays">Education Days:</label>
    <input 
      type="number" 
      id="educationDays"
      min="0" 
      max="100"
      value={schedule.EducationDays}
      on:input={(e) => handleInputChange('EducationDays', e)}
    />
  </div>

  <div class="form-group">
    <label for="sickDays">Sick Days:</label>
    <input 
      type="number" 
      id="sickDays"
      min="0" 
      max="50"
      value={schedule.SickDays}
      on:input={(e) => handleInputChange('SickDays', e)}
    />
  </div>

  <div class="schedule-summary">
    <h4>Schedule Summary</h4>
    <p><strong>Working Hours per Week:</strong> {schedule.HoursPerWeek}</p>
    <p><strong>Non-working Days:</strong> {schedule.VacationDays + schedule.PublicHolidays + schedule.EducationDays + schedule.SickDays} days/year</p>
    <p><strong>Estimated Working Days:</strong> {Math.max(0, 365 - (schedule.VacationDays + schedule.PublicHolidays + schedule.EducationDays + schedule.SickDays))} days/year</p>
  </div>
</div>

<style>
  .schedule-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .form-group label {
    font-weight: 500;
    color: #333;
  }

  .form-group input {
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
  }

  .form-group input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
  }

  .schedule-summary {
    margin-top: 1rem;
    padding: 1rem;
    background: #f8f9fa;
    border-radius: 4px;
    border-left: 4px solid #667eea;
  }

  .schedule-summary h4 {
    margin: 0 0 0.5rem 0;
    color: #333;
  }

  .schedule-summary p {
    margin: 0.25rem 0;
    font-size: 0.9rem;
    color: #666;
  }
</style> 