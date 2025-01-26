package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googleai"
	"google.golang.org/api/option"
	"log"
	"os"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	ctx := context.Background()

	gak := os.Getenv("GEMINI_API_KEY")
	gp := os.Getenv("GEMINI_1.5_PRO")

	if err := googleai.Init(ctx, &googleai.Config{
		APIKey:        gak,
		ClientOptions: []option.ClientOption{},
	}); err != nil {
		log.Fatalf("Failed to initialize Google AI plugin: %v", err)
	}

	out := genkit.DefineFlow("travelPlannerFlow", func(ctx context.Context, destination string) (string, error) {
		model := googleai.Model(gp)
		if model == nil {
			return "", errors.New("travelPlannerFlow: couldn't find the Gemini travel model")
		}
		promptText := fmt.Sprintf("Create a 3-day itinerary for traveling to %s, focusing on local cuisine and sights.", destination)
		gr, err := model.Generate(ctx,
			ai.NewGenerateRequest(
				&ai.GenerationCommonConfig{
					MaxOutputTokens: 2048, // Reduced for example
					Temperature:     0.7,
					//Consider setting a specific Version or removing if default is desired
				},
				ai.NewUserTextMessage(promptText),
			),
			nil,
		)
		if err != nil {
			return "", fmt.Errorf("travelPlannerFlow: failed to generate itinerary: %w", err)
		}

		if gr == nil {
			return "", errors.New("travelPlannerFlow: generated response is nil")
		}

		itinerary := gr.Text()
		fmt.Println(itinerary)

		return itinerary, nil
	})

	s := out.Stream(ctx, "London")
	for chunk := range s {
		fmt.Println("stream: ", chunk)
	}
	// Initialize and un Genkit's flows server.
	if err := genkit.Init(ctx, nil); err != nil {
		log.Fatalf("Failed to initialize Genkit: %v", err)
	}

	log.Println("Genkit for Go - Travel PlannerFLow is up and running!")
}
