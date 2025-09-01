# Use a Node.js Alpine image for the builder stage
FROM node:22-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build
RUN npm prune --production

# Use another Node.js Alpine image for the final stage
FROM node:22-alpine
WORKDIR /app
COPY --from=builder /app/build build/
COPY --from=builder /app/node_modules node_modules/
COPY package.json .
EXPOSE 5173
ENV PORT=5173
ENV NODE_ENV=production
CMD [ "node", "build" ]


