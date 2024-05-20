from locust import HttpUser, task

class BroadcastMessages(HttpUser):
    @task
    def broadcast(self):
        id = "2"
        data = {"content": "value"}  
        self.client.post(f"/api/channels/{id}", json=data)
        
    @task
    def get_channels(self):
        self.client.get("/api/channel_details")