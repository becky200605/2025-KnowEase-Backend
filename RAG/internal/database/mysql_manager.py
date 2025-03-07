from typing import List, Dict
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, Session
from sqlalchemy.pool import QueuePool
from internal.config.config import Config
from .models import Base, Post


class MySQLManager:
    def __init__(self, config: Config = None):
        if config is None:
            config = Config()
        mysql_config = config.get_mysql_config()
        
        # 构建数据库 URL
        db_url = f"mysql+pymysql://{mysql_config['user']}:{mysql_config['password']}@{mysql_config['host']}:{mysql_config['port']}/{mysql_config['database']}"
        
        # 创建引擎
        self.engine = create_engine(
            db_url,
            poolclass=QueuePool,
            pool_size=5,
            max_overflow=10,
            pool_timeout=30,
            pool_recycle=3600
        )
        
        # 创建会话工厂
        self.SessionLocal = sessionmaker(
            autocommit=False,
            autoflush=False,
            bind=self.engine
        )
        
        print("🔧 初始化数据库连接池...")
        
    def init_db(self):
        """初始化数据库表"""
        print("🏗️  创建数据库表...")
        Base.metadata.create_all(bind=self.engine)
        print("✅ 数据库表创建完成")
        
    def get_session(self) -> Session:
        """获取数据库会话"""
        return self.SessionLocal()
        
    def get_new_posts(self, from_id: int = 0) -> List[Dict]:
        """获取新帖子"""
        print(f"🔍 获取 ID > {from_id} 的新帖子...")
        
        session = self.get_session()
        try:
            posts = session.query(Post)\
                .filter(Post.id > from_id)\
                .order_by(Post.id.asc())\
                .all()
                
            result = [post.to_dict() for post in posts]
            print(f"📝 获取到 {len(result)} 条新帖子")
            print(Post.question)
            return result
            
        except Exception as e:
            print(f"❌ 获取帖子失败: {str(e)}")
            raise
        finally:
            session.close()