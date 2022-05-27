
## Mock Services

- [ ] Ability to create a REST mock service
	- [ ] Ability to add mock endpoints for REST mock services
- [ ] Ability to create a GraphQL mock service 🎯 
	- [ ] Ability to add a mock for a query 🎯
	- [ ] Ability to create a mock for a mutation 🎯
- [ ] Ability to configure global response headers
- [ ] Ability to create scenarios of type 'HTTP Response'
- [ ] Ability to create scenarios of type 'Proxy'
	- [ ] Ability to configure proxy request headers and proxy response headers in proxy scenarios
- [ ] Partial Proxying: If enabled, Ability to configure an upstream endpoint for mock service, such that requests to any non-implemented endpoints in the mock service will be proxied to the upstream
	- [ ] Ability to configure proxy request headers and proxy response headers in partial proxy mode
- [ ] Should be able activate a random scenario, from the configured scenarios.
- [ ] Ability to create scenarios of type 'Network Condition'
	- [ ] Ability to create scenario 'Socket Hang up'
- [ ] Ability to configure a response delay
	- [ ] Ability to send a response body with randomly generated data using faker.js
	- [ ] Ability to send a dynamic response depending on the content of the request, using JSONata mapping
- [ ] Ability to update mock endpoints
	- [ ] Ability to delete a scenario
- [ ] Ability to delete mock endpoints
- [ ] Ability to view mock service configuration
- [ ] Ability to view all mock services
- [ ] Ability to configure CORS for a mock service

## Workspaces
- [ ] Ability to create a workspace
- [ ] Ability to update a workspace
- [ ] Ability to list workspaces
- [ ] Ability to delete a workspace
- [ ] Ability to view request/response log for the workspace

## History
- [ ] Ability to view request/response log for workspace
- [ ] Ability to view request/response log for each endpoint in workspace
- [ ] Ability to clear activity log for workspace
- [ ] Ability to export activity log for workspace

## Importing
- [ ] Ability to import an OpenAPI Spec and auto generate the mock service
- [ ] Ability to import a Switchboard dump file

## Exporting
- [ ] Ability to export a mock service config as a Switchboard dump file

## Configuration API
- [ ] Ability to create an API client that generates a Client ID/Client Secret key pair
