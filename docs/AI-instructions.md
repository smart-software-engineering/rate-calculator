# AI Integration Instructions

## Overview

This document describes the technical implementation for integrating AI-driven rate calculations into the Rate Calculator system.

## Data Structure

The system uses rich JSON templates for cities and routes with embedded AI hints and metadata. The data follows a consistent schema with:
- `metadata` section (location info, coordinates, etc.)
- `referenceData` section (sources and averages)
- `businessContext` section (areas, hubs, and accommodation)

## Technical Implementation

### 1. Core Data Structures

```go
// Request structure for AI calculations
type RateCalculationRequest struct {
    HomeLocation      string       `json:"homeLocation"`
    WorkLocation      string       `json:"workLocation"`
    TravelFrequency   string       `json:"travelFrequency"`
    Context          RateContext   `json:"context"`
}

type RateContext struct {
    CityData        map[string]CityInfo    `json:"cityData"`
    RouteData       map[string]RouteInfo   `json:"routeData"`
    Preferences     UserPreferences        `json:"preferences"`
}

// Interface for AI service
type AIRateService interface {
    CalculateOptimalRate(req RateCalculationRequest) (*RateRecommendation, error)
    GenerateExplanation(recommendation *RateRecommendation) string
}
```

### 2. OpenAI Implementation

```go
type OpenAIRateService struct {
    client    *openai.Client
    prompts   map[string]string
}

func (s *OpenAIRateService) CalculateOptimalRate(req RateCalculationRequest) (*RateRecommendation, error) {
    // 1. Build context from JSON data
    context := buildPromptContext(req)
    
    // 2. Create chat completion request
    messages := []openai.ChatCompletionMessage{
        {
            Role:    "system",
            Content: SystemPrompt,
        },
        {
            Role:    "user",
            Content: fmt.Sprintf("Calculate optimal rate for travel from %s to %s with context: %s",
                req.HomeLocation, req.WorkLocation, context),
        },
    }
    
    // 3. Make API call with structured output
    resp, err := s.client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model:       "gpt-4",
            Messages:    messages,
            Temperature: 0.2,
            Functions: []openai.FunctionDefinition{
                {
                    Name: "calculate_rate",
                    Parameters: map[string]interface{}{
                        "type": "object",
                        "properties": {
                            "baseRate": {"type": "number"},
                            "travelPremium": {"type": "number"},
                            "accommodationCosts": {"type": "number"},
                            "reasoning": {"type": "string"},
                        },
                    },
                },
            },
        },
    )
}
```

### 3. System Prompt

```go
const SystemPrompt = `You are an AI rate calculator that helps determine optimal rates for freelancers.
You have access to detailed city and route information, including:
- Cost of living indices
- Travel routes and times
- Accommodation costs
- Local business contexts
Your task is to analyze this data and provide rate recommendations that account for:
1. Base rate adjusted for local cost of living
2. Travel time and costs
3. Accommodation and daily expenses
4. Work impact factors (timezone, energy levels, etc)`
```

### 4. Response Handling

```go
type RateRecommendation struct {
    BaseRate            float64   `json:"baseRate"`
    TravelPremium      float64   `json:"travelPremium"`
    AccommodationCosts float64   `json:"accommodationCosts"`
    TotalRate          float64   `json:"totalRate"`
    Reasoning          string    `json:"reasoning"`
    Factors            []string  `json:"factors"`
}

func processAIResponse(resp openai.ChatCompletionResponse) (*RateRecommendation, error) {
    var result RateRecommendation
    if err := json.Unmarshal([]byte(resp.Choices[0].Message.FunctionCall.Arguments), &result); err != nil {
        return nil, fmt.Errorf("failed to parse AI response: %w", err)
    }
    
    if err := validateRecommendation(&result); err != nil {
        return nil, fmt.Errorf("invalid rate recommendation: %w", err)
    }
    
    return &result, nil
}
```

### 5. Context Building

```go
func buildPromptContext(req RateCalculationRequest) string {
    // Extract relevant data
    homeCity := req.Context.CityData[req.HomeLocation]
    workCity := req.Context.CityData[req.WorkLocation]
    route := req.Context.RouteData[fmt.Sprintf("%s-%s", req.HomeLocation, req.WorkLocation)]
    
    // Structure the context
    context := map[string]interface{}{
        "home_location": map[string]interface{}{
            "cost_index": homeCity.Metadata.CostIndex,
            "economic_zone": homeCity.Metadata.EconomicZone,
            "typical_costs": homeCity.ReferenceData.Averages,
        },
        "work_location": map[string]interface{}{
            "cost_index": workCity.Metadata.CostIndex,
            "business_areas": workCity.BusinessContext.MainBusinessAreas,
            "accommodation": workCity.BusinessContext.TypicalAccommodation,
        },
        "route": map[string]interface{}{
            "travel_time": route.RecommendedTransport.Primary.Duration,
            "typical_costs": route.RecommendedTransport.Primary.Typical.Cost,
            "workable_time": route.WorkableTime,
        },
    }
    
    return json.MarshalToString(context)
}
```

### 6. Validation

```go
func validateRecommendation(rec *RateRecommendation) error {
    if rec.BaseRate <= 0 {
        return errors.New("base rate must be positive")
    }
    if rec.TravelPremium < 0 {
        return errors.New("travel premium cannot be negative")
    }
    if rec.Reasoning == "" {
        return errors.New("recommendation must include reasoning")
    }
    return nil
}
```

## Integration Notes

- Place the AI service in the `internal/calculator` directory
- Call the AI service after collecting user input but before displaying rate recommendations
- Use structured function calls to ensure consistent, well-formatted responses
- Keep temperature low (0.2) for consistent numerical outputs
- Include validation at both request and response stages
- Store system prompts in a configuration file for easy updates

## Error Handling

The implementation includes multiple layers of error handling:
1. Input validation before API calls
2. API error handling
3. Response parsing validation
4. Business logic validation
5. Rate recommendation validation

## Example Usage

The AI service should be called from your rate calculation endpoint:

```go
func (h *Handler) CalculateRate(w http.ResponseWriter, r *http.Request) {
    var req RateCalculationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    recommendation, err := h.aiService.CalculateOptimalRate(req)
    if err != nil {
        http.Error(w, "Failed to calculate rate", http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(recommendation)
}
