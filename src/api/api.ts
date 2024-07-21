// Packages
import cors from "@fastify/cors";
import ratelimit from "@fastify/rate-limit";
import swagger from "@fastify/swagger";
import ui from "@fastify/swagger-ui";
import Fastify, { FastifyInstance } from "fastify";
import autoload from "@fastify/autoload";
import config from "../config.js";
import { success } from "../logger.js";

// Config
const { name, description, version, port } = config;

// Middleware
const app: FastifyInstance = Fastify({
	logger: false,
});

app.register(cors, {
	origin: "*",
	allowedHeaders: [
		"Authorization",
		"Content-Type",
		"Content-Disposition",
		"Content-Length",
	],
	methods: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"],
	credentials: true,
	optionsSuccessStatus: 200,
	preflight: true,
	strictPreflight: false,
});

app.register(swagger, {
	openapi: {
		info: {
			title: name.charAt(0).toUpperCase() + name.slice(1),
			description: description,
			version: version,
		},
		tags: [
			{
				name: "platforms",
				description:
					"API endpoints responsible for accessing platform information.",
			},
			{
				name: "plugins",
				description:
					"API endpoints responsible for accessing plugin information.",
			},
		],
	},
	hideUntagged: true,
});

app.register(ui, {
	routePrefix: "/docs",
	uiConfig: {
		docExpansion: "full",
		deepLinking: true,
	},
	uiHooks: {
		onRequest: (request, reply, next) => next(),
		preHandler: (request, reply, next) => next(),
	},
	staticCSP: true,
	transformStaticCSP: (header) => header,
	transformSpecification: (swaggerObject, request, reply) => swaggerObject,
	transformSpecificationClone: true,
});

app.register(ratelimit, {
	global: true,
	max: 50,
	timeWindow: 1000,
});

app.addHook("preHandler", (req, res, done) => {
	res.header("Access-Control-Allow-Origin", "*");
	res.header("Access-Control-Allow-Headers", "*");
	res.header("Access-Control-Allow-Methods", "*");
	res.header("Access-Control-Allow-Credentials", "true");

	done();
});

app.register(autoload, {
	dir: "./dist/api/endpoints",
});

// Swagger
app.ready(() => app.swagger());

// Start Server
app.listen({ port: Number(port), host: "0.0.0.0" }, (err) => {
	if (err) throw err;
	else
		success(
			"API",
			`Server has started and is listening on http://0.0.0.0:${String(
				port
			)}/ 🚀🌐`
		);
});
