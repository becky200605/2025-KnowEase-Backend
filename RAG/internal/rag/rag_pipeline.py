from haystack import Pipeline
from haystack.utils import Secret
from haystack_integrations.document_stores.elasticsearch import ElasticsearchDocumentStore
from haystack_integrations.components.retrievers.elasticsearch import ElasticsearchEmbeddingRetriever
from haystack.components.generators import OpenAIGenerator
from haystack.components.builders.prompt_builder import PromptBuilder
from haystack.components.embedders import SentenceTransformersTextEmbedder, SentenceTransformersDocumentEmbedder
from haystack.document_stores.types import DuplicatePolicy
from internal.config.config import Config

class RAGPipeline:
    def __init__(self, config: Config = None):
        if config is None:
            config = Config()
            
        es_config = config.get_elasticsearch_config()
        llm_config = config.get_llm_config()
        
        # 定义索引映射
        custom_mapping = {
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
                        #"author_id": {"type": "keyword"},
                        #"author_name":{"type":"text"},
                        #"created_at": {"type": "date"},
                        #"tag": {"type": "text"},
                        "post_id":{"type":"text"}
                    }
                }
            }
        }
        
        # Initialize document store with embedding configuration
        self.document_store = ElasticsearchDocumentStore(
            hosts=es_config["hosts"][0],
            index="qqq",
            embedding_similarity_function="cosine",  # 使用余弦相似度
            custom_mapping=custom_mapping
        )
        self.llm_config = llm_config
        self.model_name = llm_config["embedding_model"]
        self.pipeline = self._create_pipeline()
        
    def _create_pipeline(self):
        prompt_template = """
        你叫小知，是一个专业的社区问答助手。请基于以下文档内容，分析用户所提问题和文档中问题的相似度，返回相似度并用专业且友好的语气回答用户的问题。
        如果文档中没有相关内容，请诚实的告诉用户未找到相关问答，并提示可以去问答社区进行发帖询问。
        如果文档中有类似文档内容，请在具体回答中返回该问题的post_id，以便做出下一步跳转操作,并请根据文档中信息回答用户的问题。
        所有答案不为空的问答文档都不进行相似度计算，只生成回答。
        额外信息：华中师范大学别名可以是ccnu,也可以是嘻嘻恩尤等等。

        文档内容:
        {% for doc in documents %}
              - 问题和答案: {{ doc.content }}
              - 帖子ID:{{ doc.meta.post_id }}
              - 帖子答案：{{doc.meta.answer}}
        {% endfor %}


        问题: {{question}}

        回答结构:
        - **相似度**: 输出问题和文档中内容的相似度，取值范围为 0 到 1。
        - **具体回答**: 基于文档内容，提供详细的答案。
        - **帖子ID**:基于文档内容查找相似的帖子，如果文档中存在该问题而未查询到答案，请查找并返回该问题的post_id。


        请以 JSON 格式返回结果，格式如下：
        {
          "similarity": <相似度值>,
          "answer": "<具体回答>",
          "postids":"<帖子ID>"
        }
        """
        
        # Initialize embedders
        self.text_embedder = SentenceTransformersTextEmbedder(model=self.model_name)
        self.text_embedder.warm_up()

        # Initialize retriever with embedding configuration
        self.retriever = ElasticsearchEmbeddingRetriever(
            document_store=self.document_store,
            top_k=5
        )
        
        self.prompt_builder = PromptBuilder(template=prompt_template)
        
        self.llm = OpenAIGenerator(
            api_key=Secret.from_token(self.llm_config["api_key"]),
            model=self.llm_config["model"],
            api_base_url=self.llm_config["base_url"],
        )
        
        pipeline = Pipeline()
        pipeline.add_component("text_embedder", self.text_embedder)
        pipeline.add_component("retriever", self.retriever)
        pipeline.add_component("prompt_builder", self.prompt_builder)
        pipeline.add_component("llm", self.llm)

        # Connect components
        pipeline.connect("text_embedder.embedding", "retriever.query_embedding")
        pipeline.connect("retriever", "prompt_builder.documents")
        pipeline.connect("prompt_builder", "llm")
        
        return pipeline
        
    def query(self, question: str):
        print(f"❓ 处理查询: {question}")
        print("🔄 运行 RAG pipeline...")
        results = self.pipeline.run(
            {
                "text_embedder": {"text": question},
                "prompt_builder": {"question": question},
            }
        )
        print("✅ RAG pipeline 运行完成")
        return results
        
    def query_related_documents(self, question: str):
        print(f"❓ 处理查询: {question}")
        print("🔄 运行 RAG pipeline...")
        embedding = self.text_embedder.run(question)

        results = self.retriever.run(embedding['embedding'])
        return results
    
