## Idea Validation Tool

This project is designed to help users validate a user idea by gathering external sentiment from Reddit and leveraging LLM capabilities for structured planning.

### Project Overview

The tool facilitates the process of taking a user idea and systematically validating it by:
1. **Planning**: Generating a structured plan for the idea.
2. **Research**: Checking public sentiment on Reddit to gauge community interest.
3. **Execution**: Using integrated LLMs and tools to assist in the validation process.

### Project Structure

The project is primarily written in Go and is structured around several components:

*   **`providers/`**: Contains implementations for external services, such as `gemini.go` and `ollama.go`.
*   **`tools/`**: Houses the logic for interacting with external sources, including `reddit.go` for Reddit scraping/querying and `search.go` for general searching.
*   **`ui.go`**: Contains the core logic for the user interface or main application flow.
*   **`prompts/`**: Stores predefined prompts used to guide the LLM in generating plans and research strategies, such as `growth_stratergy.txt` and `generate-idea.txt`.
*   **`info.txt`**: Contains general information about the project and workflow.

### How to Use

1.  **Setup**: Ensure your environment variables (e.g., in `.env`) are configured for any necessary API keys or credentials.
2.  **Idea Generation**: Use the available prompts to guide the initial phase of idea shaping.
3.  **Validation Workflow**: Follow the steps defined in the project's workflow (detailed in `info.txt`) to execute the validation process, which involves:
    *   Using the research tools (e.g., Reddit) to gather external data.
    *   Utilizing the LLM providers to synthesize this information into an actionable plan.
4.  **Execution**: The `ui.go` component orchestrates the calls between the tools, providers, and prompts to guide you through the validation of your idea.