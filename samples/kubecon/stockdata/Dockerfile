# Use official Python 3.9 slim image as the base
FROM python:3.9-slim

# Set working directory in the container
WORKDIR /app

# Copy the application code to the container
COPY stock_server.py .

# Install dependencies
RUN pip install --no-cache-dir flask yfinance pandas

# Expose the port the app will run on
EXPOSE 5001

# Command to run the application
CMD ["python", "stock_server.py"]