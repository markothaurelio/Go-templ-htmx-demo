# Go HTMX Clean Architecture

A structured web application (with basic news article functionality) built with **Go**, **HTMX**, **Templ**, and
**PostgreSQL**, organized using clean architecture principles.

This project demonstrates layered backend design, authentication
handling, and database integration in a production-style structure.

------------------------------------------------------------------------

## Tech Stack

-   Go\
-   HTMX\
-   Templ\
-   PostgreSQL\
-   Docker

------------------------------------------------------------------------

## Project Structure

    handlers/        HTTP routing and request handling
    services/        Business logic
    repositories/    Database access layer
    models/          Domain models
    middleware/      Authentication & cross-cutting logic
    templates/       UI components (Templ)
    postgres_db/     Database setup

------------------------------------------------------------------------

## Features

-   Layered architecture
-   Authentication flow
-   PostgreSQL integration
-   Server-rendered UI with HTMX
-   Docker support

## Purpose

Designed to showcase backend engineering structure and clean separation
of concerns using modern Go web tooling.
