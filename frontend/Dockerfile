FROM node:18-alpine

WORKDIR /app

# Copy package.json and package-lock.json (if available)
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application
COPY . .

# Expose Vite's default dev port
EXPOSE 5173

# Run Vite dev server, bound to 0.0.0.0 so it's accessible outside the container
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
