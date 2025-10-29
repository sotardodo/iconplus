# ============================
# Stage 1: Build frontend
# ============================
FROM node:18 AS build
WORKDIR /app

# Copy dependency files and install
COPY package*.json ./
RUN npm install --legacy-peer-deps

# Copy source code and build
COPY . .
RUN npm run build

# ============================
# Stage 2: Run static server
# ============================
FROM node:18-alpine
WORKDIR /app
ENV NODE_ENV=production

# Install static server
RUN npm install -g serve --only=production

# Copy built app from previous stage
COPY --from=build /app/build ./build


# Expose port 3000
EXPOSE 3000

CMD ["serve", "-s", "build", "-l", "tcp://0.0.0.0:3000"]
