import os
import yaml
from typing import Dict

class Config:
    def __init__(self, config_path: str = "config/config.yaml"):
        with open(config_path, 'r') as f:
            self.config = yaml.safe_load(f)
            
    def get_mysql_config(self) -> Dict:
        return self.config['mysql']
        
    def get_elasticsearch_config(self) -> Dict:
        return self.config['elasticsearch']
        
    def get_llm_config(self) -> Dict:
        return self.config['llm'] 
