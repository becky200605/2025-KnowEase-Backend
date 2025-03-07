import os
import grpc
from concurrent import futures
import time
from typing import List, Dict

from ..database.mysql_manager import MySQLManager
from ..elasticsearch.es_manager import ElasticsearchManager
from ..rag.rag_pipeline import RAGPipeline
from ..config.config import Config

# 导入生成的 protobuf 代码
from .generated import rag_service_pb2 as rag_pb2
from .generated import rag_service_pb2_grpc as rag_pb2_grpc

class RAGServicer(rag_pb2_grpc.RAGServiceServicer):
    def __init__(self):
        print("🚀 初始化 RAG 服务...")
        self.config = Config()
        print("⚙️  配置加载完成")
        
        print("🔄 初始化 MySQL 管理器...")
        self.mysql_manager = MySQLManager(config=self.config)
        print("✅ MySQL 管理器初始化完成")
        
        print("🔄 初始化 Elasticsearch 管理器...")
        self.es_manager = ElasticsearchManager(config=self.config)
        print("✅ Elasticsearch 管理器初始化完成")
        
        print("🔄 初始化 RAG Pipeline...")
        self.rag_pipeline = RAGPipeline(config=self.config)
        print("✅ RAG Pipeline 初始化完成")
        
    def Search(self, request, context):
        print(f"\n🔍 收到搜索请求: {request.query}")
        try:
            print("📚 开始检索相关文档...")
            # 先获取相关文档
            retrieval_results = self.rag_pipeline.query_related_documents(request.query)
            documents = retrieval_results['documents']
            print(f"📝 找到 {len(documents)} 个相关文档")
            print(documents)
            # 生成回答
            print("🤖 生成回答...")
            results = self.rag_pipeline.query(request.query)
            answer = results["llm"]["replies"][0]
            print("💡 生成回答完成")
            
            # 转换文档格式
            response_documents = []
            for doc in documents[:8]: 
                response_documents.append(rag_pb2.Document(
                    id=doc.meta.get("id",0),
                    post_id=doc.meta.get("post_id", ""),
                    question=doc.meta.get("question",""),
                    answer=doc.meta.get("answer",""),
                    #author_name=doc.meta.get("author_name", ""),
                    #author_id=doc.meta.get("author_id", ""),
                    #tag=doc.meta.get("tag",""),
                    #created_at=doc.meta.get("created_at", "")
                ))
            
            print("✅ 搜索请求处理完成")
            return rag_pb2.SearchResponse(
                documents=response_documents,
                answer=answer
            )
        except Exception as e:
            print(f"❌ 搜索请求处理失败: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return rag_pb2.SearchResponse()
            
    def SyncData(self, request, context):
        print(f"\n🔄 收到同步请求: from_id={request.from_id}")
        try:
            print("📥 从 MySQL 获取新数据...")
            posts = self.mysql_manager.get_new_posts(request.from_id)
            print(f"📊 获取到 {len(posts)} 条新数据")
            
            if posts:
                print("📤 开始更新 Elasticsearch...")
                self.es_manager.index_posts(posts)
                print("✅ Elasticsearch 更新完成")
                
            print("✅ 同步请求处理完成")
            return rag_pb2.SyncDataResponse(
                synced_count=len(posts)
            )
        except Exception as e:
            print(f"❌ 同步请求处理失败: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return rag_pb2.SyncDataResponse()

def serve():
    print("\n🚀 === 启动 RAG 服务 ===")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    rag_pb2_grpc.add_RAGServiceServicer_to_server(RAGServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("✨ 服务启动完成，监听端口: 50051")
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        print("\n🛑 收到终止信号，正在关闭服务...")
        server.stop(0)
        print("👋 服务已关闭")

if __name__ == '__main__':
    serve() 