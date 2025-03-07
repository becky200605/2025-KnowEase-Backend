from internal.database.models import Base
from internal.database.mysql_manager import MySQLManager
from internal.config.config import Config

def create_tables():
    print("🚀 开始初始化数据库...")
    config = Config()
    manager = MySQLManager(config)
    manager.init_db()
    print("✅ 数据库初始化完成")

if __name__ == '__main__':
    create_tables() 