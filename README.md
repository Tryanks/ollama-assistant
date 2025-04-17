# Ollama Assistant

A Go-based web service for interacting with AI models.

## Setup

### Environment Variables

This project requires environment variables to be set up before running. Follow these steps:

1. Copy the `.env.example` file to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Open the `.env` file and update the values:
   ```
   API_BASE_URL=https://api.openai.com
   API_KEY=<KEY>
   ```
   
   Replace `<KEY>` with your actual API key.

## Roadmap

The following is the development roadmap for this project:

- [ ] **Multi API Provider Support**: Integration with various AI service providers including OpenAI, Anthropic, local Ollama, etc.
- [ ] **System Service Implementation**: Running as a stable background system service
- [ ] **Auto-Start with System**: Adding capability to automatically start with the system
- [ ] **Memory Optimization & Performance Improvement**: Trimming memory allocation and making the program faster and more responsive
- [ ] **GUI Configuration Tool**: Developing a graphical interface for easy management and configuration

## Running the Application

After setting up the environment variables, you can run the application:

```bash
go run .
```

The server will start on port 11434.

## API Endpoints

- `GET /`: Check if the service is running
- `GET /api/tags`: List available models
- `POST /api/chat`: Chat completion endpoint

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.