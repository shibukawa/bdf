#include <blend2d/blend2d.h>
#include <fstream>
#include <iostream>
#include <vector>
#include <chrono>
#include <cstring>
struct Op { BLPath path; BLRgba32 color; bool even; double stroke; };
int main(int argc,char**argv) {
  std::ifstream in(argv[1]); int w,h,n; in>>w>>h>>n; std::vector<Op> ops;
  for(int i=0;i<n;i++){unsigned r,g,b,a;std::string rule;int ns;double stroke;in>>r>>g>>b>>a>>rule>>stroke>>ns;Op op;op.color=BLRgba32(r,g,b,a);op.even=rule=="evenodd";op.stroke=stroke;
    for(int j=0;j<ns;j++){int k;double x,y,a,b,c,d;in>>k;switch(k){case 0:in>>x>>y;op.path.move_to(x,y);break;case 1:in>>x>>y;op.path.line_to(x,y);break;case 2:in>>a>>b>>x>>y;op.path.quad_to(a,b,x,y);break;case 3:in>>a>>b>>c>>d>>x>>y;op.path.cubic_to(a,b,c,d,x,y);break;case 4:op.path.close();break;default:return 2;}}
    ops.push_back(std::move(op));}
  BLImage image(w,h,BL_FORMAT_PRGB32);BLContextCreateInfo info{};info.thread_count=argc>3?std::stoi(argv[3]):0;BLContext ctx(image,info);volatile unsigned checksum=0;
  auto render=[&](bool copy){ctx.fill_all(BLRgba32(0xFFFFFFFF));for(auto&op:ops){ctx.set_fill_rule(op.even?BL_FILL_RULE_EVEN_ODD:BL_FILL_RULE_NON_ZERO);if(op.stroke>0){ctx.set_stroke_width(op.stroke);ctx.set_stroke_miter_limit(10);ctx.stroke_path(op.path,op.color);}else ctx.fill_path(op.path,op.color);}ctx.flush(BL_CONTEXT_FLUSH_SYNC);if(copy){std::vector<unsigned char> owned(w*h*4);BLImageData data;image.get_data(&data);for(int y=0;y<h;y++)std::memcpy(owned.data()+y*w*4,(char*)data.pixel_data+y*data.stride,w*4);asm volatile("" : : "g"(owned.data()) : "memory");checksum=owned[0];}};
  auto start=std::chrono::steady_clock::now();render(true);double cold=std::chrono::duration<double,std::nano>(std::chrono::steady_clock::now()-start).count();
  std::cout<<"{\"backend\":\"blend2d-cpu\",\"threads\":"<<info.thread_count<<",\"first_render_ns\":"<<cold<<",\"results\":{";
  for(int copy=0;copy<2;copy++){if(copy)std::cout<<",";std::cout<<"\""<<(copy?"owned-output":"render-only")<<"\":{\"samples_ns\":[";for(int sample=0;sample<5;sample++){auto t=std::chrono::steady_clock::now();int count=0;do{render(copy);count++;}while(std::chrono::steady_clock::now()-t<std::chrono::milliseconds(300));if(sample)std::cout<<",";std::cout<<std::chrono::duration<double,std::nano>(std::chrono::steady_clock::now()-t).count()/count;}std::cout<<"]}";}
  std::cout<<"}}\n";ctx.end();return image.write_to_file(argv[2])!=BL_SUCCESS;
}
