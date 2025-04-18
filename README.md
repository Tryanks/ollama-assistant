# Ollama Assistant

A Go-based web service that provides an Ollama-compatible API interface while using OpenAI's API under the hood. This allows you to use Ollama clients with OpenAI models.

## Features

- **Ollama-Compatible API**: Implements Ollama's API endpoints for seamless integration with existing Ollama clients
- **OpenAI Backend**: Uses OpenAI's official Go client to communicate with OpenAI's API
- **Multi-Provider Support**: Connect to multiple OpenAI-compatible API providers simultaneously
  - Automatically routes requests to the appropriate provider based on the model name
  - Supports different base URLs and API keys for each provider
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
   # Legacy configuration (still supported)
   API_BASE_URL=https://api.openai.com
   API_KEY=<KEY>

   # New configuration format for multiple providers
   # Format: OPENAI_PROVIDERS=url1, key1; url2, key2;...
   OPENAI_PROVIDERS=

   # Server configuration
   HOST_SERVE=0.0.0.0
   PORT_SERVE=11434
   ```

   Replace `<KEY>` with your actual OpenAI API key.

   - `API_BASE_URL` and `API_KEY`: Legacy configuration for a single OpenAI provider
   - `OPENAI_PROVIDERS`: New configuration format for multiple OpenAI providers (takes precedence if set)
     - Format: `OPENAI_PROVIDERS=url1, key1; url2, key2;...`
     - Example: `OPENAI_PROVIDERS=https://api.openai.com, sk-xxx; https://api.another-provider.com, sk-yyy`
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
- `GET /api/refresh-models`: Refresh the cached model list and return the updated list
- `POST /api/chat`: Chat completion endpoint that accepts Ollama-format requests and returns Ollama-format responses

## Performance Optimizations

- **Cached Model List**: The application scans the model list only once at startup and caches it for future requests, improving performance by reducing API calls to OpenAI
- **Manual Refresh**: Use the `/api/refresh-models` endpoint to manually refresh the cached model list when needed

## Usage with Ollama Clients

You can use any Ollama client with this service. Simply point the client to this service's URL instead of the official Ollama service.

## Docker Support

### Using the Docker Image

The application is available as a Docker image from GitHub Container Registry:

```bash
# Pull the latest image
docker pull ghcr.io/tryanks/ollama-assistant:latest

# Run the container
docker run -d \
  -p 11434:11434 \
  -e API_KEY=your_openai_api_key \
  -e API_BASE_URL=https://api.openai.com \
  ghcr.io/tryanks/ollama-assistant:latest
```

### Building the Docker Image Locally

You can build the Docker image locally using the provided Dockerfile:

```bash
docker build -t ollama-assistant .
docker run -d -p 11434:11434 -e API_KEY=your_openai_api_key ollama-assistant
```

### Automated Builds with GitHub Actions

This repository includes a GitHub Actions workflow that automatically builds and publishes the Docker image to GitHub Container Registry (ghcr.io) when:
- Changes are pushed to the main branch
- A new tag is created (prefixed with 'v')
- The workflow is manually triggered

To set up automated builds:

1. Create a GitHub Personal Access Token (PAT) with `read:packages` and `write:packages` scopes
2. Add the token as a repository secret named `GITHUB_PERSONAL_ACCESS_TOKEN`
3. Push changes to the main branch or create a new tag to trigger the workflow

The workflow will build and tag the image with:
- The semantic version (for tags)
- The branch name (for branch pushes)
- The commit SHA
- 'latest' tag for the most recent build

## Roadmap

The following is the development roadmap for this project:

- [x] **Basic OpenAI Integration**: Core functionality with OpenAI API
- [x] **Environment Configuration**: Support for configurable host and port settings
- [x] **Multi OpenAI Provider Support**: Support for multiple OpenAI-compatible API providers
- [x] **Docker Support**: Containerized deployment option
- [ ] **Other API Provider Support**: Integration with Anthropic Claude and other non-OpenAI providers
- [ ] **Local Ollama Fallback**: Option to use local Ollama when API is unavailable
- [ ] **Performance Optimization**: Improved response handling and memory management
- [ ] **Simple Web UI**: Basic web interface for testing and configuration

## Knowledge

> The code in this repository was (almost entirely) written by the JetBrains AI Assistant itself through its ollama integration.  
> Except for the initial phase of basic chat API translation.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
