from locust import HttpUser, task

class BroadcastMessages(HttpUser):
    @task
    def broadcast_2(self):
        id = "2"
        data = {"content": "value: two"}  
        self.client.post(f"/api/v1/channels/{id}", json=data)

    # @task
    # def broadcast_3(self):
    #     id = "3"
    #     data = {"content": "value: three"}  
    #     self.client.post(f"/api/channels/{id}", json=data)
        
    # @task
    # def broadcast_4(self):
    #     id = "4"
    #     data = {"content": "value: four"}  
    #     self.client.post(f"/api/channels/{id}", json=data)
        
    # @task
    # def get_channels(self):
    #     self.client.get("/api/channel/all")