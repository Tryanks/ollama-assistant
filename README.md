# Ollama Assistant

A Go-based web service that provides an Ollama-compatible API interface while using OpenAI's API under the hood. This allows you to use Ollama clients with OpenAI models.

## Features

- **Ollama-Compatible API**: Implements Ollama's API endpoints for seamless integration with existing Ollama clients
- **OpenAI Backend**: Uses OpenAI's official Go client to communicate with OpenAI's API
- **Advanced AI Features Support**:
  - Streaming responses
  - Tool calling (function calling)
  - Structured outputs (JSON mode)
- **Model Filtering**: Automatically filters out non-chat models like DALL-E, embeddings, TTS, and Whisper

## Setup

### Environment Variables

This project requires environment variables to be set up before running. Follow these steps:

1. Create a `.env` file in the project root:
   ```bash
   touch .env
   ```

2. Open the `.env` file and add the following values:
   ```
   API_BASE_URL=https://api.openai.com
   API_KEY=<KEY>
   HOST_SERVE=0.0.0.0
   PORT_SERVE=11434
   ```

   Replace `<KEY>` with your actual OpenAI API key.

   - `API_BASE_URL`: The base URL for the OpenAI API
   - `API_KEY`: Your OpenAI API key
   - `HOST_SERVE`: The host address to bind the server to (default: 0.0.0.0)
   - `PORT_SERVE`: The port to run the server on (default: 11434)

## Running the Application

After setting up the environment variables, you can run the application:

```bash
go run .
```

By default, the server will start on port 11434 (the same port used by the official Ollama server), but you can change this using the PORT_SERVE environment variable.

## API Endpoints

- `GET /`: Check if the service is running
- `GET /api/tags`: List available OpenAI chat models (filtered to exclude non-chat models)
- `POST /api/chat`: Chat completion endpoint that accepts Ollama-format requests and returns Ollama-format responses

## Usage with Ollama Clients

You can use any Ollama client with this service. Simply point the client to this service's URL instead of the official Ollama server.

## Roadmap

The following is the development roadmap for this project:

- [x] **Basic OpenAI Integration**: Core functionality with OpenAI API
- [x] **Environment Configuration**: Support for configurable host and port settings
- [ ] **Multi API Provider Support**: Integration with Anthropic Claude and other providers
- [ ] **Local Ollama Fallback**: Option to use local Ollama when API is unavailable
- [ ] **Performance Optimization**: Improved response handling and memory management
- [ ] **Docker Support**: Containerized deployment option
- [ ] **Simple Web UI**: Basic web interface for testing and configuration

## Knowledge

> The code in this repository was (almost entirely) written by the JetBrains AI Assistant itself through its ollama integration.  
> Except for the initial phase of basic chat API translation.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
