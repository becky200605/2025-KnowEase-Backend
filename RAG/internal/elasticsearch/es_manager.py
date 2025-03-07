from haystack_integrations.document_stores.elasticsearch import ElasticsearchDocumentStore
from haystack import Document
from haystack.document_stores.types import DuplicatePolicy
from haystack.components.embedders import SentenceTransformersDocumentEmbedder
from typing import List, Dict, Set
from internal.config.config import Config
from elasticsearch import Elasticsearch, NotFoundError

class ElasticsearchManager:
    def __init__(self, config: Config = None):
        if config is None:
            config = Config()
        es_config = config.get_elasticsearch_config()
        

        self.client = Elasticsearch(
            hosts=es_config["hosts"][0],
            verify_certs=False
        )
        
        # 定义索引映射
        self.custom_mapping = {
            "properties": {
                "content": {"type": "text"},
                "embedding": {
                    "type": "dense_vector",
                    "dims": 1024,  
                    "similarity": "cosine"
                },
                "meta": {
                    "properties": {
                        #"question":{"type": "text"},
                        "answer":{"type":"text"},
                        "id": {"type": "long"},
                        "post_id":{"type":"text"},
                        #"author_id": {"type": "keyword"},
                        #"author_name":{"type":"text"},
                        #"created_at": {"type": "date"},
                        #"tag": {"type": "text"}
                    }
                }
            }
        }
        
        self.index_name = "qqq"
        self._ensure_index_exists()
        
        self.document_store = ElasticsearchDocumentStore(
            hosts=es_config["hosts"][0],
            index=self.index_name,
            embedding_similarity_function="cosine",  # 使用余弦相似度
            custom_mapping=self.custom_mapping
        )
        
        self.model_name = "BAAI/bge-large-zh-v1.5"
        self.document_embedder = SentenceTransformersDocumentEmbedder(model=self.model_name)
        self.document_embedder.warm_up()
        
        # 用于缓存已存在的文档ID
        self._existing_ids: Set[int] = set()
        
        
    def _ensure_index_exists(self):
        """确保索引存在，如果不存在则创建"""
        try:
            if not self.client.indices.exists(index=self.index_name):
                print(f"📦 创建索引 {self.index_name}...")
                self.client.indices.create(
                    index=self.index_name,
                    mappings=self.custom_mapping
                )
                print("✅ 索引创建成功")
        except Exception as e:
            print(f"⚠️ 检查/创建索引时出错: {str(e)}")
            raise
            
    def _check_document_exists(self, doc_id: int) -> bool:
        """检查文档是否存在"""
        try:
            # 首先确保索引存在
            if not self.client.indices.exists(index=self.index_name):
                return False
                
            query = {
                "query": {
                    "term": {
                        "meta.id": doc_id
                    }
                }
            }
            result = self.client.search(
                index=self.index_name,
                body=query,
                size=1
            )
            return result["hits"]["total"]["value"] > 0
        except NotFoundError:
            return False
        except Exception as e:
            print(f"⚠️ 检查文档存在性时出错: {str(e)}")
            return False
            
    def _filter_new_posts(self, posts: List[Dict]) -> List[Dict]:
        """过滤出新的帖子"""
        new_posts = []
        for post in posts:
            if not self._check_document_exists(post["id"]):
                new_posts.append(post)
            else:
                print(f"⏩ 跳过已存在的帖子: ID={post['id']}")
                
        skipped = len(posts) - len(new_posts)
        if skipped > 0:
            print(f"⏩ 共跳过 {skipped} 条已存在的帖子")
        return new_posts
        
    def index_posts(self, posts: List[Dict]):
        """索引新帖子"""
        try:
            # 确保索引存在
            self._ensure_index_exists()
            
            # 过滤已存在的帖子
            new_posts = self._filter_new_posts(posts)
            if not new_posts:
                print("ℹ️ 没有新的帖子需要索引")
                return
                
            print(f"📥 索引 {len(new_posts)} 条新文档...")
            documents = []
            for post in new_posts:
                doc = Document(
                    content=f"Question: {post['question']}\n\nAnswer: {post['answer']}",
                    meta={
                        #"question":post["question"],
                        "answer":post["answer"],
                        "id": post["id"],
                        "post_id":post["post_id"],
                        #"author_id": post["author_id"],
                        #"author_name": post["author_name"],
                        #"created_at": post["created_at"].isoformat(),
                        #"tag": post["tag"]
                    }
                )
                documents.append(doc)
            #print(documents)
                
            print("📥 生成文档嵌入...")
            documents_with_embeddings = self.document_embedder.run(documents)
            # 打印生成的向量值
            #for doc in documents_with_embeddings["documents"]:
            #    print(f"文档 ID: {doc.meta['id']}, 嵌入向量: {doc.embedding}")
                
            print("💾 写入文档到 Elasticsearch...")
            self.document_store.write_documents(
                documents_with_embeddings["documents"],
                policy=DuplicatePolicy.SKIP
            )
            #print(documents_with_embeddings)
            print("✅ 文档索引完成")
        except Exception as e:
            print(f"❌ 写入文档失败: {str(e)}")
            raise
