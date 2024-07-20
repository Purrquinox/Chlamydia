// Packages
import fs from "node:fs";
import path from "path";
import cors from "@fastify/cors";
import ratelimit from "@fastify/rate-limit";
import swagger from "@fastify/swagger";
import ui from "@fastify/swagger-ui";
import Fastify, { FastifyInstance } from "fastify";
import { name, description, version } from "../../package.json";
import { success } from "../logger";

// Middleware
const PORT: Number = 50805;
const app: FastifyInstance = Fastify({
	logger: true,
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
	swagger: {
		info: {
			title: name.charAt(0).toUpperCase() + name.slice(1),
			description: description,
			version: version,
		},
		host: `localhost:${PORT}`,
		schemes: ["http"],
		consumes: ["application/json"],
		produces: ["application/json"],
		tags: [
			{
				name: "platforms",
				description:
					"API endpoints responsible for accessing platform information.",
			},
		],
		securityDefinitions: {},
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
		onRequest: (request, reply, next) => {
			next();
		},
		preHandler: (request, reply, next) => {
			next();
		},
	},
	staticCSP: true,
	transformStaticCSP: (header) => header,
	transformSpecification: (swaggerObject, request, reply) => {
		return swaggerObject;
	},
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

// API Endpoints Map
const getFilesInDirectory = (dir: string) => {
	let files: string[] = [];
	const filesInDir = fs.readdirSync(dir);

	for (const file of filesInDir) {
		const filePath = path.join(dir, file);
		const stat = fs.statSync(filePath);

		if (stat.isDirectory())
			files = files.concat(getFilesInDirectory(filePath));
		else files.push(filePath);
	}

	return files;
};

// API Endpoints
const apiEndpointsFiles = getFilesInDirectory(
	"./dist/src/api/endpoints"
).filter((file) => file.endsWith(".js"));

for (const file of apiEndpointsFiles) {
	import(`../../../${file}`)
		.then(async (module) => {
			await app.route(module.default.default);
		})
		.catch((error) => {
			console.error(`Error importing ${file}: ${error}`);
		});
}

setTimeout(() => {
	// Swagger
	app.ready(() => {
		app.swagger();
	});

	// Start Server
	app.listen({ port: Number(PORT) }, (err) => {
		if (err) throw err;
		else
			success(
				"API",
				`Server has started and is listening on http://localhost:${String(
					PORT
				)}/ 🚀🌐`
			);
	});
}, 8000);
