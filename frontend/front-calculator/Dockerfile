# Calculadora web — imagen de producción.
# Contexto de build: frontend/front-calculator/
#   docker build -t front-calculator ./frontend/front-calculator
#
# VITE_API_BASE_URL se fija en build (Vite la incrusta). Para otro backend:
#   docker build --build-arg VITE_API_BASE_URL=http://api:8080 -t front-calculator ./frontend/front-calculator
# Nota: el navegador llama a esa URL directamente, debe ser alcanzable desde el navegador.

FROM node:22-alpine AS build

WORKDIR /app
RUN corepack enable

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .

ARG VITE_API_BASE_URL=http://localhost:8080
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
RUN pnpm build

FROM nginx:1.27-alpine AS runtime

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
