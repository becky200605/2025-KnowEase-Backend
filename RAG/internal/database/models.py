from sqlalchemy import Column, Integer, String, Text, DateTime, func
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()

class Post(Base):
    __tablename__ = 'QAs'
    
    id = Column(Integer, primary_key=True, autoincrement=True)
    post_id = Column(String(100),nullable=False)
    question = Column(Text, nullable=False)
    answer = Column(Text, nullable=True)
    #author_name = Column(String(100), nullable=False)
    #author_id = Column(String(100), nullable=False)
    #author_url = Column(String(255), nullable=False)
    #tag = Column(String(255), nullable=False)
    #created_at = Column(DateTime, default=func.now())
    #updated_at = Column(DateTime, default=func.now(), onupdate=func.now())
    
    def to_dict(self):
        return {
            'id': self.id,
            'post_id':self.post_id,
            'question':self.question,
            'answer':self.answer,
            #'tag':self.tag,
            #'author_name':self.author_name,
            #'author_id': self.author_id,
            #'author_url':self.author_url,
            #'created_at': self.created_at,
            #'updated_at': self.updated_at
        }